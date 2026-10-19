//go:build integration

package engagelab

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

type qaList struct{ V []any }

func (v *qaList) MarshalJSON() ([]byte, error) { return json.Marshal(v.V) }
func qarr(v ...any) any {
	if v == nil {
		v = []any{}
	}
	return &qaList{v}
}
func qlistOf(v any) []any {
	switch x := v.(type) {
	case *qaList:
		if x != nil {
			return x.V
		}
	case []any:
		return x
	case []string:
		r := []any{}
		for _, s := range x {
			r = append(r, s)
		}
		return r
	}
	return []any{}
}
func qarg(v []any, i int) any {
	if i < len(v) {
		return v[i]
	}
	return nil
}
func qm(v ...any) map[string]any {
	r := map[string]any{}
	for i := 0; i+1 < len(v); i += 2 {
		r[qtxt(v[i])] = v[i+1]
	}
	return r
}
func qg(v any, k ...any) any {
	for _, key := range k {
		switch x := v.(type) {
		case map[string]any:
			if value, present := x[qtxt(key)]; present {
				v = value
			} else {
				v = nil
				for k, value := range x {
					if fullKey(k) == fullKey(qtxt(key)) {
						v = value
						break
					}
				}
			}
		default:
			a := qlistOf(v)
			i := int(qnum(key))
			if i < 0 || i >= len(a) {
				return nil
			}
			v = a[i]
		}
	}
	return v
}
func qdrop(v, k any) any     { delete(v.(map[string]any), qtxt(k)); return nil }
func qput(v, k, val any) any { v.(map[string]any)[qtxt(k)] = val; return nil }
func qappend(v, x any) any   { a := v.(*qaList); a.V = append(a.V, x); return nil }
func qtxt(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprint(v)
}
func qnum(v any) int64     { n, _ := strconv.ParseInt(qtxt(v), 10, 64); return n }
func qcat(a, b any) string { return qtxt(a) + qtxt(b) }
func qsize(v any) int64 {
	switch a := v.(type) {
	case string:
		return int64(len([]rune(a)))
	case map[string]any:
		return int64(len(a))
	}
	return int64(len(qlistOf(v)))
}
func qyes(v any) bool { b, ok := v.(bool); return ok && b }
func qno(v any) bool  { b, ok := v.(bool); return ok && !b }
func qnonempty(v any) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case string, map[string]any, []any, *qaList, []string:
		return qsize(v) > 0
	}
	return true
}
func qeq(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	if _, ok := a.(bool); ok {
		return reflect.DeepEqual(a, b)
	}
	if _, ok := b.(bool); ok {
		return false
	}
	na, oka := new(big.Rat).SetString(qtxt(a))
	nb, okb := new(big.Rat).SetString(qtxt(b))
	if oka && okb {
		return na.Cmp(nb) == 0
	}
	if ma, ok := a.(map[string]any); ok {
		mb, yes := b.(map[string]any)
		if !yes || len(ma) != len(mb) {
			return false
		}
		for k, v := range ma {
			w, p := mb[k]
			if !p || !qeq(v, w) {
				return false
			}
		}
		return true
	}
	if len(qlistOf(a)) > 0 || len(qlistOf(b)) > 0 {
		aa, bb := qlistOf(a), qlistOf(b)
		if len(aa) != len(bb) {
			return false
		}
		for i := range aa {
			if !qeq(aa[i], bb[i]) {
				return false
			}
		}
		return true
	}
	return reflect.DeepEqual(a, b)
}
func qboth(v ...any) bool {
	for _, x := range v {
		if !qyes(x) {
			return false
		}
	}
	return true
}
func qeither(v ...any) bool {
	for _, x := range v {
		if qyes(x) {
			return true
		}
	}
	return false
}
func qincludes(a, v any) bool {
	if x, ok := a.(map[string]any); ok {
		_, p := x[qtxt(v)]
		return p
	}
	if x, ok := a.(string); ok {
		return strings.Contains(x, qtxt(v))
	}
	for _, x := range qlistOf(a) {
		if qeq(x, v) {
			return true
		}
	}
	return false
}
func qsubset(a, b any) bool {
	for _, x := range qlistOf(a) {
		if !qincludes(b, x) {
			return false
		}
	}
	return true
}
func qsets(a, b any) bool { return qsubset(a, b) && qsubset(b, a) }
func qsame(a, b any) bool {
	if want, ok := b.(map[string]any); ok {
		actual, ok := a.(map[string]any)
		if !ok {
			return false
		}
		keys := map[string]any{}
		for k, v := range actual {
			keys[fullKey(k)] = v
		}
		for k, v := range want {
			if k == "mobile" || k == "target_converted" || k == "fileMd5" || k == "file_md5" {
				continue
			}
			got, p := keys[fullKey(k)]
			if !p || !qsame(got, v) {
				return false
			}
		}
		return true
	}
	switch b.(type) {
	case []any, *qaList:
		aa, bb := qlistOf(a), qlistOf(b)
		if len(aa) != len(bb) {
			return false
		}
		for i := range bb {
			if !qsame(aa[i], bb[i]) {
				return false
			}
		}
		return true
	}
	return qeq(a, b)
}
func qkeys(v any) any {
	r := qarr()
	if x, ok := v.(map[string]any); ok {
		for k := range x {
			qappend(r, k)
		}
	}
	return r
}
func qids(v, key any) any {
	r := qarr()
	for _, x := range qlistOf(v) {
		qappend(r, qg(x, key))
	}
	return r
}
func qjoined(v any) string {
	r := []string{}
	for _, x := range qlistOf(v) {
		r = append(r, qtxt(x))
	}
	return strings.Join(r, ",")
}
func qquery(v any) string {
	q := url.Values{}
	for k, val := range v.(map[string]any) {
		q.Set(k, qtxt(val))
	}
	return q.Encode()
}
func qspan(v any) any {
	r := qarr()
	for i := int64(0); i < qnum(v); i++ {
		qappend(r, i)
	}
	return r
}
func qtake(v, n any) any {
	a := qlistOf(v)
	k := int(qnum(n))
	if k > len(a) {
		k = len(a)
	}
	return qarr(a[:k]...)
}
func qwithout(v, x any) any {
	r := qarr()
	for _, z := range qlistOf(v) {
		if !qeq(z, x) {
			qappend(r, z)
		}
	}
	return r
}
func qintersects(a, b any) bool {
	for _, v := range qlistOf(a) {
		if qincludes(b, v) {
			return true
		}
	}
	return false
}
func qge(a, b any) bool { return qnum(a) >= qnum(b) }
func qgt(a, b any) bool { return qnum(a) > qnum(b) }
func qleft(a, n any) string {
	s := []rune(qtxt(a))
	k := int(qnum(n))
	if k > len(s) {
		k = len(s)
	}
	return string(s[:k])
}
func qrepeat(a, n any) string { return strings.Repeat(qtxt(a), int(qnum(n))) }
func qadd(a, b any) int64     { return qnum(a) + qnum(b) }
func qmaxi(a, b any) int64 {
	if qnum(a) > qnum(b) {
		return qnum(a)
	}
	return qnum(b)
}
func qcoalesce(v ...any) any {
	for _, x := range v {
		if qnonempty(x) {
			return x
		}
	}
	return nil
}
func qfirstOrSelf(v any) any {
	switch v.(type) {
	case []any, *qaList:
		return qg(v, 0)
	}
	return v
}
func qisInteger(v any) bool {
	switch v.(type) {
	case int, int64, int32, uint64, json.Number:
		return true
	}
	return false
}
func qhasKeysShallow(v, keys any) bool {
	for _, k := range qlistOf(keys) {
		m, ok := v.(map[string]any)
		if !ok {
			return false
		}
		if _, p := m[qtxt(k)]; !p {
			return false
		}
	}
	return true
}
func qhasKeysDeep(v, keys any) bool {
	var has func(any, string) bool
	has = func(v any, key string) bool {
		if m, ok := v.(map[string]any); ok {
			if _, p := m[key]; p {
				return true
			}
			for _, x := range m {
				if has(x, key) {
					return true
				}
			}
		}
		for _, x := range qlistOf(v) {
			if has(x, key) {
				return true
			}
		}
		return false
	}
	for _, k := range qlistOf(keys) {
		if !has(v, qtxt(k)) {
			return false
		}
	}
	return true
}
func qhmosNonzero(v any) bool {
	if m, ok := v.(map[string]any); ok {
		for k, x := range m {
			if strings.EqualFold(k, "hmos") {
				if n, ok := x.(map[string]any); ok {
					for _, value := range n {
						if qnum(value) > 0 {
							return true
						}
					}
				} else if qnum(x) > 0 {
					return true
				}
			}
			if qhmosNonzero(x) {
				return true
			}
		}
	}
	for _, x := range qlistOf(v) {
		if qhmosNonzero(x) {
			return true
		}
	}
	return false
}
func qparse(b []byte) any {
	var v any
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	if d.Decode(&v) != nil {
		return nil
	}
	return v
}
func qclone(v any) any { b, _ := json.Marshal(v); return qparse(b) }

// Preserve typed zero/false and 64-bit values, including fields hidden from JSON marshal.
func qtree(v any) any {
	if v == nil {
		return nil
	}
	if g, ok := v.(*GroupPushResult); ok {
		return qm("group_msgid", g.GroupMsgID, "successes", qtree(g.Successes), "errors", qtree(g.Errors))
	}
	r := reflect.ValueOf(v)
	for r.Kind() == reflect.Ptr || r.Kind() == reflect.Interface {
		if r.IsNil() {
			return nil
		}
		r = r.Elem()
	}
	switch r.Kind() {
	case reflect.Struct:
		m := qm()
		t := r.Type()
		for i := 0; i < r.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" {
				continue
			}
			k := strings.Split(f.Tag.Get("json"), ",")[0]
			if k == "-" {
				continue
			}
			if k == "" {
				k = f.Name
			}
			m[k] = qtree(r.Field(i).Interface())
		}
		return m
	case reflect.Map:
		m := qm()
		it := r.MapRange()
		for it.Next() {
			m[qtxt(it.Key().Interface())] = qtree(it.Value().Interface())
		}
		return m
	case reflect.Slice, reflect.Array:
		a := qarr()
		for i := 0; i < r.Len(); i++ {
			qappend(a, qtree(r.Index(i).Interface()))
		}
		return a
	default:
		return r.Interface()
	}
}

type qaBlocked struct{ reason string }
type fullQA struct {
	cfg      nativeConfig
	envs     map[string]string
	prefix   string
	actions  *fullActions
	last     map[string]any
	network  string
	results  map[string]any
	owned    [][2]string
	aliases  map[string]string
	fixtures []string
	states   map[string]any
	http     *http.Client
}

func newFullQA() *fullQA {
	e := nativeEnv()
	for _, k := range []string{"ENGAGELAB_LIFECYCLE_MESSAGE_ID", "ENGAGELAB_LIFECYCLE_REGISTRATION_ID"} {
		e[k] = strings.TrimSpace(e[k])
	}
	if (e["ENGAGELAB_LIFECYCLE_MESSAGE_ID"] != "") != (e["ENGAGELAB_LIFECYCLE_REGISTRATION_ID"] != "") {
		panic(qaBlocked{"Provide both lifecycle message ID and registration ID"})
	}
	cfg, err := nativeLoad(e)
	if err != nil {
		panic(qaBlocked{err.Error()})
	}
	c := &fullQA{cfg: cfg, envs: e, prefix: cfg.Run, results: qm(), aliases: map[string]string{}, states: qm()}
	c.http = &http.Client{Timeout: cfg.Timeout, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	c.actions = newFullActions(cfg.Base, cfg.Key, cfg.Secret, e["ENGAGELAB_GROUP_KEY"], e["ENGAGELAB_GROUP_SECRET"], cfg.Timeout, c)
	for _, r := range strings.Split(e["ENGAGELAB_FIXTURE_REGISTRATION_IDS"], ",") {
		if r != "" {
			c.fixtures = append(c.fixtures, r)
		}
	}
	return c
}
func (c *fullQA) RoundTrip(r *http.Request) (*http.Response, error) {
	time.Sleep(c.cfg.Interval)
	raw := []byte{}
	if r.Body != nil {
		raw, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(raw))
	}
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewReader(body))
	query := qm()
	for k, v := range r.URL.Query() {
		query[k] = qtree(v)
	}
	c.last = qm("status", resp.StatusCode, "raw", qparse(body), "request", qparse(raw), "path", r.URL.Path, "query", query, "multipart_valid", strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data; boundary=") && bytes.Contains(raw, []byte(`name="language"`)) && bytes.Contains(raw, []byte(`name="file"`)) && bytes.Contains(raw, c.voiceBytes()))
	return resp, err
}
func (c *fullQA) sdk(op, args any) any {
	c.last = nil
	c.network = ""
	raw, _ := json.Marshal(args)
	var params []json.RawMessage
	if json.Unmarshal(raw, &params) != nil {
		panic("native argument serialization failed")
	}
	value, err := c.actions.call(qtxt(op), params)
	out := qm()
	for k, v := range c.last {
		out[k] = v
	}
	if err == nil {
		out["ok"] = true
		out["data"] = qtree(value)
		return out
	}
	out["ok"] = false
	if ae, ok := err.(*ApiError); ok {
		out["error_code"] = ae.ErrorBody.Code
		out["error_message"] = ae.ErrorBody.Message
	} else if c.last == nil {
		c.network = "SDK transport unavailable"
		out["network"] = c.network
	} else {
		panic("native SDK response could not decode")
	}
	return out
}
func (c *fullQA) request(key, secret, method, path string, body any) any {
	time.Sleep(c.cfg.Interval)
	var raw []byte
	if body != nil {
		raw, _ = json.Marshal(body)
	}
	req, err := http.NewRequest(method, c.cfg.Base+path, bytes.NewReader(raw))
	if err != nil {
		panic("invalid control URL")
	}
	req.SetBasicAuth(key, secret)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		panic(qaBlocked{"QA control transport unavailable"})
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(qaBlocked{"QA control body unavailable"})
	}
	return qm("status", resp.StatusCode, "data", qparse(b))
}
func (c *fullQA) rest(method, path any, body ...any) any {
	return c.request(c.cfg.Key, c.cfg.Secret, qtxt(method), qtxt(path), qarg(body, 0))
}
func (c *fullQA) require(ok, why any) any {
	if !qyes(ok) {
		panic(qaBlocked{qtxt(why)})
	}
	return nil
}
func (c *fullQA) write() any {
	if c.cfg.Mode != "isolated_write" {
		panic(qaBlocked{"Requires isolated_write"})
	}
	return nil
}
func (c *fullQA) target() any { return c.cfg.Env }
func (c *fullQA) need(key any) string {
	v := strings.TrimSpace(c.envs[qtxt(key)])
	if v == "" {
		panic(qaBlocked{"Missing " + qtxt(key)})
	}
	return v
}
func (c *fullQA) env(key, def any) any {
	if v := c.envs[qtxt(key)]; v != "" {
		return v
	}
	return def
}
func (c *fullQA) envNumber(key, def any) any { return qnum(c.env(key, def)) }
func (c *fullQA) unique() any                { return fmt.Sprintf("%x", time.Now().UnixNano()) }
func (c *fullQA) now() any                   { return time.Now().Unix() }
func (c *fullQA) dateTime(ts any) any {
	return time.Unix(qnum(ts), 0).In(time.FixedZone("UTC+8", 28800)).Format("2006-01-02 15:04:05")
}
func (c *fullQA) date(ts any) any {
	return time.Unix(qnum(ts), 0).In(time.FixedZone("UTC+8", 28800)).Format("2006-01-02")
}
func (c *fullQA) sleepUntil(ts any) any {
	if d := time.Until(time.Unix(qnum(ts), 0)); d > 0 {
		time.Sleep(d)
	}
	return nil
}
func (c *fullQA) state(k, v any) any {
	if v != nil {
		c.states[qtxt(k)] = v
	}
	return c.states[qtxt(k)]
}
func (c *fullQA) own(kind, value any) any {
	c.owned = append(c.owned, [2]string{qtxt(kind), qtxt(value)})
	dir := os.Getenv("ENGAGELAB_FULL_RESULTS_DIR")
	if dir == "" {
		dir = ".test-results/full"
	}
	dir = filepath.Join(dir, "resources")
	if os.MkdirAll(dir, 0700) != nil {
		panic("resource journal directory unavailable")
	}
	b, _ := json.Marshal(qm("run_id", c.prefix, "owned", c.owned))
	if os.WriteFile(filepath.Join(dir, c.prefix+".json"), b, 0600) != nil {
		panic("resource journal unavailable")
	}
	return nil
}
func (c *fullQA) ownedValues(kind, max any) any {
	r := qarr()
	for _, v := range c.owned {
		if v[0] == qtxt(kind) && qsize(r) < qnum(max) {
			qappend(r, v[1])
		}
	}
	return r
}
func (c *fullQA) rememberAlias(rid any) any {
	r := qtxt(rid)
	for _, v := range c.owned {
		if v[0] == "device" && v[1] == r {
			return nil
		}
	}
	if _, ok := c.aliases[r]; !ok {
		c.aliases[r] = qtxt(qg(c.device(r), "data", "alias"))
	}
	return nil
}
func (c *fullQA) fixtureRID(index any) any {
	i := int(qnum(index))
	if i < len(c.fixtures) {
		return c.fixtures[i]
	}
	return nil
}
func (c *fullQA) addFixture(rid any) any { c.fixtures = append(c.fixtures, qtxt(rid)); return nil }
func (c *fullQA) waitFor(kind, value any) any {
	if !qnonempty(value) {
		return false
	}
	end := time.Now().Add(c.cfg.Wait)
	for {
		ready := false
		switch qtxt(kind) {
		case "plans":
			r := c.listPlans(qg(value, "prefix"), 101, 1, nil)
			ready = qeq(qg(r, "data", "total"), qg(value, "count"))
		case "device_absent":
			ready = qeq(qg(c.device(value), "data", "error", "code"), 27002)
		case "schedule":
			r := c.rest("GET", "/v4/schedules/"+qtxt(value)+"/msg-ids", nil)
			ready = qnonempty(qcoalesce(qg(r, "data", "msg_ids"), qg(r, "data", "msgIds")))
		case "messages":
			r := c.rest("GET", "/v4/status/detail?"+qquery(qm("message_ids", qjoined(value))), nil)
			ready = true
			for _, mid := range qlistOf(value) {
				ready = ready && qnonempty(qg(r, "data", mid))
			}
		}
		if ready {
			return true
		}
		if time.Now().After(end) {
			return false
		}
		time.Sleep(c.cfg.Interval)
	}
}
func (c *fullQA) record(sid any, checks map[string]any, ready bool, why string) {
	status := "pass"
	for _, v := range checks {
		if !qyes(v) {
			status = "no_pass"
		}
	}
	reason := ""
	if !ready {
		status = "block"
		reason = why
	} else if c.network != "" {
		status = "block"
		reason = c.network
	}
	c.results[qtxt(sid)] = qm("status", status, "checks", checks, "reason", reason)
}
func (c *fullQA) clean() {
	failed := []string{}
	safe := func(label string, fn func() bool) {
		defer func() {
			if recover() != nil {
				failed = append(failed, label)
			}
		}()
		if !fn() {
			failed = append(failed, label)
		}
	}
	plans := qlistOf(c.ownedValues("plan", 10000))
	for i := 0; i < len(plans); i += 10 {
		end := i + 10
		if end > len(plans) {
			end = len(plans)
		}
		chunk := plans[i:end]
		safe("plan_delete", func() bool { c.rest("DELETE", "/v4/push_plan/batch/"+qjoined(qarr(chunk...)), nil); return true })
	}
	if len(plans) > 0 {
		safe("plan_readback", func() bool { return qyes(c.waitFor("plans", qm("prefix", c.prefix, "count", 0))) })
	}
	for i := len(c.owned) - 1; i >= 0; i-- {
		kind, val := c.owned[i][0], c.owned[i][1]
		if kind == "plan" {
			continue
		}
		safe(kind+"_readback", func() bool {
			switch kind {
			case "alias":
				p := "/v4/aliases/" + val + "?platform=android"
				c.rest("DELETE", p, nil)
				return !qnonempty(qg(c.rest("GET", p, nil), "data", "registration_ids"))
			case "tag":
				c.rest("DELETE", "/v4/tags/"+val+"?platform=android", nil)
				return !qincludes(qg(c.rest("GET", "/v4/tags/", nil), "data", "tags"), val)
			case "device":
				c.rest("DELETE", "/v4/devices/"+val, nil)
				return qyes(c.waitFor("device_absent", val))
			case "schedule":
				c.rest("DELETE", "/v4/schedules/"+val, nil)
				return qeq(qg(c.rest("GET", "/v4/schedules/"+val, nil), "status"), 404)
			case "voice":
				c.rest("DELETE", "/v4/voices/"+val, nil)
				return qeq(qg(c.rest("GET", "/v4/voices/"+val, nil), "status"), 404)
			}
			return false
		})
	}
	for rid, alias := range c.aliases {
		r, a := rid, alias
		safe("fixture_alias_restore", func() bool {
			c.rest("POST", "/v4/devices/"+r, qm("alias", a))
			return qeq(qtxt(qg(c.device(r), "data", "alias")), a)
		})
	}
	if wav, ok := c.states["wav"]; ok {
		_ = os.Remove(qtxt(wav))
	}
	for _, result := range c.results {
		r := result.(map[string]any)
		qg(r, "checks").(map[string]any)["owned_resource_cleanup"] = len(failed) == 0
		if len(failed) > 0 && qeq(r["status"], "pass") {
			r["status"] = "no_pass"
			r["reason"] = "Cleanup not verified: " + strings.Join(failed, ",")
		}
	}
}
func fullExecute(t *testing.T, cache map[string]map[string]any, group, sid string, ids []string, flow func(*fullQA)) {
	t.Helper()
	if _, ok := cache[group]; !ok {
		results := qm()
		var c *fullQA
		func() {
			defer func() {
				if err := recover(); err != nil {
					status, reason := "no_pass", "Native workflow exception"
					if b, ok := err.(qaBlocked); ok {
						status, reason = "block", b.reason
					}
					for _, id := range ids {
						if c == nil || c.results[id] == nil {
							results[id] = qm("status", status, "checks", qm(), "reason", reason)
						}
					}
				}
				if c != nil {
					c.clean()
					for k, v := range c.results {
						results[k] = v
					}
				}
			}()
			c = newFullQA()
			flow(c)
		}()
		cache[group] = results
	}
	r := cache[group][sid].(map[string]any)
	checks := []string{}
	for k, v := range qg(r, "checks").(map[string]any) {
		if qyes(v) {
			checks = append(checks, k)
		} else {
			checks = append(checks, "FAILED:"+k)
		}
	}
	fullWriteResult(t, sid, "integration", qtxt(r["status"]), qtxt(r["reason"]), checks)
	if !qeq(r["status"], "pass") {
		t.Error(qtxt(r["status"]) + ": " + qtxt(r["reason"]) + " " + strings.Join(checks, ","))
	}
}
func (c *fullQA) authControl() any {
	old := c.actions
	c.actions = newFullActions(c.cfg.Base, "000000000000000000000000", "invalid-sdk-secret", "", "", c.cfg.Timeout, c)
	r := c.sdk("vip", qarr())
	c.actions = old
	raw := c.request("000000000000000000000000", "invalid-sdk-secret", "GET", "/v4/app/vip/status", nil)
	return qm("verified", c.errorsEqual(r, raw))
}
func (c *fullQA) apnsConfig() any {
	token := c.need("ENGAGELAB_APNS_TOKEN")
	prod := c.need("ENGAGELAB_APNS_PRODUCTION")
	c.need("ENGAGELAB_BUNDLE_ID")
	c.require(regexp.MustCompile(`^[0-9a-fA-F]{64}$`).MatchString(token) && (prod == "true" || prod == "false"), "Real 64-hex APNs token and explicit channel required")
	return qm("token", token, "production", prod == "true")
}
func (c *fullQA) groupConfig() any {
	c.need("ENGAGELAB_GROUP_KEY")
	c.need("ENGAGELAB_GROUP_SECRET")
	apps := strings.Split(c.need("ENGAGELAB_GROUP_APP_KEYS"), ",")
	return qm("apps", qtree(apps), "alias", c.need("ENGAGELAB_GROUP_TARGET_ALIAS"))
}
func (c *fullQA) memberAliasBinding(app, alias any) any {
	secrets := qparse([]byte(c.need("ENGAGELAB_GROUP_APP_SECRETS")))
	secret := qg(secrets, app)
	c.require(qnonempty(secret), "Group member control credential missing")
	return c.request(qtxt(app), qtxt(secret), "GET", "/v4/aliases/"+url.PathEscape(qtxt(alias))+"?platform=android", nil)
}
func (c *fullQA) imageResource() any {
	u := c.need("ENGAGELAB_OPPO_IMAGE_URL")
	r, err := c.http.Get(u)
	if err != nil {
		panic(qaBlocked{"Image download unavailable"})
	}
	defer r.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(r.Body, 1048577))
	valid := r.StatusCode == 200 && len(b) > 24 && len(b) < 1048576 && bytes.Equal(b[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) && binary.BigEndian.Uint32(b[16:]) == 984 && binary.BigEndian.Uint32(b[20:]) == 369
	c.require(valid, "Reachable 984x369 PNG below 1MiB required")
	return qm("url", u, "valid", valid)
}
func (c *fullQA) removeFile(path any) any {
	if os.Remove(qtxt(path)) != nil {
		panic("WAV cleanup failed")
	}
	return nil
}
func (c *fullQA) adb(args ...string) string {
	path := qtxt(c.env("ENGAGELAB_ADB_PATH", "adb"))
	if len(args) > 0 && args[0] == "shell" {
		quoted := []string{}
		for _, a := range args[1:] {
			quoted = append(quoted, "'"+strings.ReplaceAll(a, "'", "'\"'\"'")+"'")
		}
		args = []string{"shell", strings.Join(quoted, " ")}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, append([]string{"-s", c.need("ENGAGELAB_ADB_SERIAL")}, args...)...)
	b, err := cmd.Output()
	if err != nil {
		panic(qaBlocked{"Declared ADB device or command unavailable"})
	}
	return string(b)
}
func (c *fullQA) notifications(marker string) []string {
	pkg := c.need("ENGAGELAB_ANDROID_PACKAGE")
	out := []string{}
	for _, k := range strings.Split(c.adb("shell", "cmd", "notification", "list"), "\n") {
		k = strings.TrimSpace(k)
		if strings.Contains(k, "|"+pkg+"|") {
			record := c.adb("shell", "cmd", "notification", "get", k)
			if strings.Contains(record, marker) {
				out = append(out, k)
			}
		}
	}
	return out
}
func (c *fullQA) physicalNotifications() any {
	c.write()
	rid := c.need("ENGAGELAB_TEST_REGISTRATION_ID")
	pkg := c.need("ENGAGELAB_ANDROID_PACKAGE")
	c.require(c.need("ENGAGELAB_DEVICE_APP_KEY") == c.cfg.Key, "Physical app identity differs from SDK app")
	c.require(strings.TrimSpace(c.adb("shell", "pidof", pkg)) != "", "Physical QA Demo not running")
	c.require(qeq(qg(c.device(rid), "status"), 200), "Physical RID absent from QA app")
	info := c.adb("shell", "dumpsys", "package", pkg)
	uid := regexp.MustCompile(`(?:userId|appId)=(\d+)`).FindStringSubmatch(info)
	c.require(len(uid) == 2, "App UID unavailable")
	logs := func(marker string, ids ...any) bool {
		raw := c.adb("logcat", "-d", "--uid="+uid[1], "-v", "epoch", "-t", "2000")
		for _, line := range strings.Split(raw, "\n") {
			if (strings.Contains(line, marker) || (len(ids) > 0 && qnonempty(ids[0]) && strings.Contains(line, qtxt(ids[0])))) && (strings.Contains(line, "onNotification") || strings.Contains(line, "UserReceiver") || strings.Contains(line, "MTNotificationBusiness")) {
				return true
			}
		}
		return false
	}
	marker := c.prefix + "_physical"
	sent := c.sdk("push", qarr(c.payload(rid, marker, nil)))
	mid := qg(sent, "data", "msg_id")
	posted := false
	if qnonempty(mid) {
		defer func() { c.rest("DELETE", "/v4/push/withdraw/"+qtxt(mid), nil) }()
		end := time.Now().Add(c.cfg.Wait)
		for time.Now().Before(end) {
			if len(c.notifications(marker)) > 0 {
				posted = true
				break
			}
			time.Sleep(time.Second)
		}
	}
	persisted := c.waitFor("messages", qarr(qtxt(mid)))
	pushOK := qboth(qeq(qg(sent, "status"), 200), posted, persisted, qeq(mid, qg(sent, "raw", "msg_id")), qsame(qg(sent, "data"), qg(sent, "raw")))
	logControl := logs(marker, mid)
	control := c.sdk("device_status", qarr(qm("registration_ids", qarr(rid))))
	valmarker := marker + "_validate"
	validation := c.payload(rid, valmarker, nil)
	qput(qg(validation, "body", "options"), "time_to_live", 10)
	validated := c.sdk("validate", qarr(validation))
	unexpected := false
	end := time.Now().Add(12 * time.Second)
	for time.Now().Before(end) {
		unexpected = unexpected || len(c.notifications(valmarker)) > 0
		time.Sleep(time.Second)
	}
	after := c.sdk("device_status", qarr(qm("registration_ids", qarr(rid))))
	validateOK := qboth(qeq(qg(validated, "status"), 200), !unexpected, qyes(qg(control, "data", 0, "online")), qyes(qg(after, "data", 0, "online")), posted, logControl, !logs(valmarker), qeq(qtxt(qg(validated, "data", "msg_id")), qtxt(qg(validated, "raw", "msg_id"))) && qeq(qtxt(qg(validated, "data", "request_id")), qtxt(qg(validated, "raw", "request_id"))))
	withdrawOK := false
	if posted && len(c.notifications(marker)) > 0 {
		withdrawn := c.sdk("withdraw", qarr(qtxt(mid)))
		end = time.Now().Add(c.cfg.Wait)
		for time.Now().Before(end) && len(c.notifications(marker)) > 0 {
			time.Sleep(time.Second)
		}
		withdrawOK = qeq(qg(withdrawn, "status"), 200) && len(c.notifications(marker)) == 0 && qeq(qg(withdrawn, "data", "msg_id"), mid) && qsame(qg(withdrawn, "data"), qg(withdrawn, "raw"))
	}
	return qm("push", pushOK, "push_ready", posted, "validate", validateOK, "validate_ready", posted && logControl, "withdraw", withdrawOK, "withdraw_ready", posted)
}

func (c *fullQA) voiceBytes() []byte {
	if b, ok := c.states["voice_bytes"].([]byte); ok {
		return b
	}
	return []byte("<no-native-voice-resource>")
}
func (c *fullQA) voiceFile() any {
	path := qtxt(c.env("ENGAGELAB_VOICE_FILE", "testdata/voice/en.zip"))
	z, err := zip.OpenReader(path)
	if err != nil {
		panic(qaBlocked{"Valid voice ZIP resource missing: ENGAGELAB_VOICE_FILE"})
	}
	defer z.Close()
	valid := len(z.File) >= 34
	for _, f := range z.File {
		valid = valid && strings.HasSuffix(f.Name, ".mp3") && f.UncompressedSize64 > 0
	}
	c.require(valid, "Expected complete nonempty MP3 voice ZIP")
	b, err := os.ReadFile(path)
	if err != nil {
		panic(qaBlocked{"Voice resource cannot be read"})
	}
	c.states["voice_bytes"] = b
	return path
}

func qdecodeJSON(text any) any { return qparse([]byte(qtxt(text))) }
