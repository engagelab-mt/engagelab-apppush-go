package engagelab

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNativeEnvironmentParameters(t *testing.T) {
	for _, v := range []struct{ target, base string }{{"qa", "http://api.qa.jpushoa.com"}, {"prod", "https://pushapi-sgp.engagelab.com"}} {
		e := map[string]string{"TEST_ENV": v.target, "ENGAGELAB_BASE_URL": v.base, "ENGAGELAB_APP_KEY": "fixture", "ENGAGELAB_MASTER_SECRET": "fixture", "ENGAGELAB_HTTP_TIMEOUT": "7"}
		c, err := nativeLoad(e)
		nativeError(t, err, "valid environment rejected")
		nativeAssert(t, c.Timeout.Seconds() == 7 && c.Mode == "read_only", "parameters not applied")
		if v.target == "qa" {
			e["TEST_ENV"] = "prod"
		} else {
			e["TEST_ENV"] = "qa"
		}
		_, err = nativeLoad(e)
		nativeAssert(t, err != nil, "environment mismatch accepted")
	}
}
func TestNativeRepresentativeContracts(t *testing.T) {
	var path, body string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.String()
		var p map[string]interface{}
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&p)
		}
		b, _ := json.Marshal(p)
		body = string(b)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/v4/push/validate":
			w.Write([]byte(`{"msg_id":"fixture-message"}`))
		case "/v4/push_plan/list":
			items := []PushPlanInfo{}
			for i := 0; i < 101; i++ {
				items = append(items, PushPlanInfo{PlanID: strings.Repeat("p", i+1)})
			}
			_ = json.NewEncoder(w).Encode(PushPlanListResult{Total: 101, PushPlanInfo: items})
		case "/v4/grouppush":
			w.Write([]byte(`{"group_msgid":"g1","app-a":{"msg_id":"m1"},"app-b":{"msg_id":"m2"}}`))
		case "/v4/devices/token/registration_id":
			w.Write([]byte(`{"results":[{"token":"fixture","registration_id":"r1","is_new":false,"code":0}]}`))
		default:
			w.WriteHeader(404)
		}
	}))
	defer server.Close()
	c := NewClient("fixture", "fixture", WithBaseURL(server.URL))
	t.Run("field_passthrough", func(t *testing.T) {
		r, e := c.Push.Validate(nativeCtx, nativePayload("r1", "contract"))
		nativeError(t, e, "validate failed")
		nativeAssert(t, r.MsgID == "fixture-message" && strings.Contains(body, `"auto_truncation":false`) && strings.Contains(body, `"badge_set_num":0`) && strings.Contains(body, `"is_fold":false`), "false/zero fields or response lost")
	})
	t.Run("pagination_101", func(t *testing.T) {
		r, e := c.Plan.List(nativeCtx, 1, 101, nil, "fixture")
		nativeError(t, e, "list failed")
		nativeAssert(t, r.Total == 101 && len(r.PushPlanInfo) == 101 && strings.Contains(path, "page_size=101"), "page 101 truncated")
	})
	t.Run("group_push", func(t *testing.T) {
		r, e := NewGroupPushClient("group", "secret", WithBaseURL(server.URL)).Send(nativeCtx, nativePayload("r1", "contract"))
		nativeError(t, e, "group failed")
		nativeAssert(t, r.GroupMsgID == "g1" && len(r.Successes) == 2 && len(r.Errors) == 0 && r.Successes["app-a"].MsgID == "m1", "dynamic group parse failed")
	})
	t.Run("apns_token", func(t *testing.T) {
		no := false
		r, e := c.Device.RegisterToken(nativeCtx, &DeviceTokenRegisterParam{Platform: "ios", Tokens: []string{"fixture"}, APNSProduction: &no})
		nativeError(t, e, "token failed")
		nativeAssert(t, len(r.Results) == 1 && r.Results[0].RegistrationID == "r1" && !r.Results[0].IsNew && strings.Contains(body, `"apns_production":false`), "sandbox false or token response lost")
	})
}
