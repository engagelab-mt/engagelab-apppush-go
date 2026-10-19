//go:build integration

package engagelab

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestIntegrationFieldPassthrough(t *testing.T) {
	w, c := nativeStart(t, "field_passthrough", false)
	p := nativePayload(w.need(t, "ENGAGELAB_TEST_REGISTRATION_ID"), w.C.Run)
	r, err := c.Push.Validate(nativeCtx, p)
	nativeError(t, err, "validate call failed")
	w.ok(t, r.MsgID != "" && r.MsgID == w.Last["msg_id"], "validate_result_matches_wire")
	b, _ := json.Marshal(w.Sent)
	s := string(b)
	w.ok(t, strings.Contains(s, `"auto_truncation":false`) && strings.Contains(s, `"badge_set_num":0`) && strings.Contains(s, `"is_fold":false`), "false_zero_preserved")
}
func TestIntegrationPagination101(t *testing.T) {
	w, c := nativeStart(t, "pagination_101", true)
	owned := []string{}
	prefix := w.C.Run
	path := "/v4/push_plan/list?" + url.Values{"page_index": {"1"}, "page_size": {"101"}, "search_description": {prefix}}.Encode()
	t.Cleanup(func() {
		for i := 0; i < len(owned); i += 10 {
			end := i + 10
			if end > len(owned) {
				end = len(owned)
			}
			if err := c.Plan.BatchDelete(nativeCtx, strings.Join(owned[i:end], ",")); err != nil {
				t.Error("plan cleanup request failed")
			}
		}
		err := w.until(func() bool {
			status, b, e := w.rest("GET", path)
			return e == nil && status == 200 && b["total"] == float64(0)
		})
		if err != nil {
			t.Error("cleanup readback failed")
		} else {
			w.Checks = append(w.Checks, "cleanup_zero_remaining")
		}
	})
	for i := 0; i < 101; i++ {
		pid := fmt.Sprintf("%s_%d", prefix, i)
		owned = append(owned, pid)
		r, e := c.Plan.CreateOrUpdate(nativeCtx, &PushPlanParam{PlanID: pid, PlanDescription: prefix})
		nativeError(t, e, "plan creation failed")
		w.ok(t, r.PlanID == pid, "created_plan_id_matches")
	}
	nativeError(t, w.until(func() bool {
		status, b, e := w.rest("GET", path)
		return e == nil && status == 200 && b["total"] == float64(101)
	}), "101-plan prerequisite failed")
	first, e := c.Plan.List(nativeCtx, 1, 100, nil, prefix)
	nativeError(t, e, "page 1 failed")
	second, e := c.Plan.List(nativeCtx, 2, 100, nil, prefix)
	nativeError(t, e, "page 2 failed")
	r, e := c.Plan.List(nativeCtx, 1, 101, nil, prefix)
	nativeError(t, e, "page 101 failed")
	w.ok(t, len(first.PushPlanInfo) == 100 && len(second.PushPlanInfo) == 1, "100_plus_1_pages")
	w.ok(t, r.Total == 101 && len(r.PushPlanInfo) == 101, "101_returned")
	want := map[string]bool{}
	for _, id := range owned {
		want[id] = true
	}
	got := map[string]bool{}
	for _, p := range r.PushPlanInfo {
		got[p.PlanID] = true
		w.ok(t, p.PlanDescription == prefix && p.CreateTime > 0, "plan_fields_preserved")
	}
	w.ok(t, reflect.DeepEqual(want, got), "owned_ids_exact")
	var wire PushPlanListResult
	b, _ := json.Marshal(w.Last)
	nativeError(t, json.Unmarshal(b, &wire), "wire parse failed")
	w.ok(t, reflect.DeepEqual(r, &wire), "parsed_fields_match_wire")
}
func TestIntegrationGroupPush(t *testing.T) {
	w, _ := nativeStart(t, "group_push", true)
	apps := strings.Split(w.need(t, "ENGAGELAB_GROUP_APP_KEYS"), ",")
	want := map[string]bool{}
	for _, a := range apps {
		want[a] = true
	}
	w.ok(t, len(want) >= 2 && !want[""], "two_explicit_members")
	alias := w.need(t, "ENGAGELAB_GROUP_TARGET_ALIAS")
	p := nativePayload("unused", w.C.Run)
	p.To = &PushTo{Alias: []string{alias}}
	ttl := 60
	p.Body.Options = &Options{TimeToLive: &ttl, ThirdPartyChannel: map[string]interface{}{"distribution_new": "mtpush"}}
	g := NewGroupPushClient(w.need(t, "ENGAGELAB_GROUP_KEY"), w.need(t, "ENGAGELAB_GROUP_SECRET"), WithBaseURL(w.C.Base), WithHTTPClient(&http.Client{Transport: w, Timeout: w.C.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}))
	r, e := g.Send(nativeCtx, p)
	nativeError(t, e, "group call failed")
	w.ok(t, r.GroupMsgID != "" && r.GroupMsgID == w.Last["group_msgid"], "group_id_matches_wire")
	got := map[string]bool{}
	for app, item := range r.Successes {
		got[app] = true
		raw, ok := w.Last[app].(map[string]interface{})
		w.ok(t, ok && item.MsgID != "" && raw["msg_id"] == item.MsgID, "per_app_id_matches_wire")
	}
	w.ok(t, reflect.DeepEqual(want, got) && len(r.Errors) == 0, "dynamic_members_all_succeeded")
}
func TestIntegrationAPNsToken(t *testing.T) {
	w, c := nativeStart(t, "apns_token", true)
	token := w.need(t, "ENGAGELAB_APNS_TOKEN")
	prod := w.need(t, "ENGAGELAB_APNS_PRODUCTION")
	_ = w.need(t, "ENGAGELAB_BUNDLE_ID")
	if !regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(token) || (prod != "true" && prod != "false") {
		w.Blocked = true
		w.Reason = "invalid APNs token/channel"
		t.Fatal("block: invalid APNs token/channel")
	}
	production := prod == "true"
	p := &DeviceTokenRegisterParam{Platform: "ios", Tokens: []string{token}, APNSProduction: &production}
	rid := ""
	created := false
	var baseline map[string]interface{}
	t.Cleanup(func() {
		if rid == "" {
			return
		}
		if created {
			_, _, e := w.rest("DELETE", "/v4/devices/"+rid)
			if e != nil {
				t.Error("token cleanup request failed")
			}
			e = w.until(func() bool {
				s, b, e := w.rest("GET", "/v4/devices/"+rid)
				detail, _ := b["error"].(map[string]interface{})
				return e == nil && s == 400 && detail["code"] == float64(27002)
			})
			if e != nil {
				t.Error("new token mapping cleanup failed")
			} else {
				w.Checks = append(w.Checks, "new_mapping_deleted")
			}
		} else if baseline != nil {
			s, b, e := w.rest("GET", "/v4/devices/"+rid)
			if e != nil || s != 200 || !reflect.DeepEqual(baseline, b) {
				t.Error("existing mapping changed")
			} else {
				w.Checks = append(w.Checks, "existing_mapping_preserved")
			}
		}
	})
	r, e := c.Device.RegisterToken(nativeCtx, p)
	nativeError(t, e, "token registration call failed")
	raw := w.Last
	entries, _ := raw["results"].([]interface{})
	if len(entries) == 1 {
		item, _ := entries[0].(map[string]interface{})
		rid, _ = item["registration_id"].(string)
		created, _ = item["is_new"].(bool)
	}
	w.ok(t, len(r.Results) == 1 && len(entries) == 1, "one_token_result")
	item := r.Results[0]
	w.ok(t, rid != "" && item.Token == token && item.RegistrationID == rid && item.Code == 0, "registered_token_fields")
	rawItem := entries[0].(map[string]interface{})
	v, ok := rawItem["is_new"].(bool)
	w.ok(t, ok && item.IsNew == v, "is_new_matches_wire")
	w.ok(t, w.Sent["apns_production"] == production, "explicit_apns_channel")
	nativeError(t, w.until(func() bool { s, _, e := w.rest("GET", "/v4/devices/"+rid); return e == nil && s == 200 }), "registered mapping not readable")
	_, baseline, e = w.rest("GET", "/v4/devices/"+rid)
	nativeError(t, e, "baseline read failed")
	again, e := c.Device.RegisterToken(nativeCtx, p)
	nativeError(t, e, "repeat token failed")
	w.ok(t, len(again.Results) == 1 && again.Results[0].RegistrationID == rid && !again.Results[0].IsNew, "same_mapping_reused")
}
