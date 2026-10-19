package engagelab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

type nativeConfig struct {
	Env, Base, Key, Secret, Mode, Run string
	Timeout, Wait, Interval           time.Duration
}

func nativeNeed(e map[string]string, k string) (string, error) {
	v := strings.TrimSpace(e[k])
	if v == "" {
		return "", fmt.Errorf("block: missing %s", k)
	}
	return v, nil
}
func nativeLoad(e map[string]string) (nativeConfig, error) {
	c := nativeConfig{Env: e["TEST_ENV"], Base: strings.TrimRight(e["ENGAGELAB_BASE_URL"], "/"), Mode: e["ENGAGELAB_MUTATION_MODE"]}
	u, err := url.Parse(c.Base)
	if err != nil {
		return c, fmt.Errorf("block: invalid API origin")
	}
	h := u.Hostname()
	qa := strings.HasSuffix(h, ".qa.jpushoa.com") || strings.HasSuffix(h, ".txqa.jpushoa.com")
	if (c.Env != "qa" && c.Env != "prod") || h == "" || u.User != nil || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || (u.Scheme != "http" && u.Scheme != "https") || (c.Env == "qa") != qa || (c.Env == "prod" && (u.Scheme != "https" || h == "localhost" || h == "127.0.0.1" || h == "::1")) {
		return c, fmt.Errorf("block: TEST_ENV and API origin mismatch")
	}
	if c.Key, err = nativeNeed(e, "ENGAGELAB_APP_KEY"); err != nil {
		return c, err
	}
	if c.Secret, err = nativeNeed(e, "ENGAGELAB_MASTER_SECRET"); err != nil {
		return c, err
	}
	if c.Mode == "" {
		c.Mode = "read_only"
	}
	if c.Mode != "read_only" && c.Mode != "isolated_write" {
		return c, fmt.Errorf("block: invalid mutation mode")
	}
	fields := []struct {
		k, d string
		dst  *time.Duration
	}{{"ENGAGELAB_HTTP_TIMEOUT", "25", &c.Timeout}, {"ENGAGELAB_WAIT_TIMEOUT", "30", &c.Wait}, {"ENGAGELAB_REQUEST_INTERVAL", "1.1", &c.Interval}}
	for _, f := range fields {
		v := e[f.k]
		if v == "" {
			v = f.d
		}
		seconds, err := strconv.ParseFloat(v, 64)
		if err != nil || seconds <= 0 || math.IsNaN(seconds) || math.IsInf(seconds, 0) || seconds > 86400 {
			return c, fmt.Errorf("block: invalid duration %s", f.k)
		}
		*f.dst = time.Duration(seconds * float64(time.Second))
	}
	prefix := e["ENGAGELAB_RESOURCE_PREFIX"]
	if prefix == "" {
		prefix = "sdknative"
	}
	if !regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,15}$`).MatchString(prefix) {
		return c, fmt.Errorf("block: invalid resource prefix")
	}
	c.Run = fmt.Sprintf("%s_%x", prefix, time.Now().UnixNano())
	return c, nil
}
func nativeEnv() map[string]string {
	e := map[string]string{}
	for _, v := range os.Environ() {
		p := strings.SplitN(v, "=", 2)
		e[p[0]] = p[1]
	}
	return e
}
func nativeAssert(t *testing.T, b bool, label string) {
	t.Helper()
	if !b {
		t.Fatal(label)
	}
}
func nativeError(t *testing.T, err error, label string) {
	t.Helper()
	if err != nil {
		t.Fatal(label)
	}
}

type nativeWire struct {
	C          nativeConfig
	Last, Sent map[string]interface{}
	Status     int
	Checks     []string
	Blocked    bool
	Reason     string
}

func (w *nativeWire) RoundTrip(r *http.Request) (*http.Response, error) {
	time.Sleep(w.C.Interval)
	if r.Body != nil {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body = io.NopCloser(bytes.NewReader(b))
		w.Sent = nil
		_ = json.Unmarshal(b, &w.Sent)
	}
	response, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return nil, err
	}
	response.Body = io.NopCloser(bytes.NewReader(b))
	w.Last = nil
	_ = json.Unmarshal(b, &w.Last)
	w.Status = response.StatusCode
	return response, nil
}
func (w *nativeWire) client() *Client {
	return NewClient(w.C.Key, w.C.Secret, WithBaseURL(w.C.Base), WithHTTPClient(&http.Client{Timeout: w.C.Timeout, Transport: w, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}))
}
func (w *nativeWire) rest(method, path string) (int, map[string]interface{}, error) {
	r, err := http.NewRequest(method, w.C.Base+path, nil)
	if err != nil {
		return 0, nil, err
	}
	r.SetBasicAuth(w.C.Key, w.C.Secret)
	c := &http.Client{Timeout: w.C.Timeout, Transport: w, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := c.Do(r)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, w.Last, nil
}
func (w *nativeWire) until(fn func() bool) error {
	end := time.Now().Add(w.C.Wait)
	for {
		if fn() {
			return nil
		}
		if time.Now().After(end) {
			return fmt.Errorf("readback deadline exceeded")
		}
		time.Sleep(w.C.Interval)
	}
}
func (w *nativeWire) ok(t *testing.T, b bool, label string) {
	t.Helper()
	nativeAssert(t, b, label)
	w.Checks = append(w.Checks, label)
}
func nativeStart(t *testing.T, scenario string, write bool) (*nativeWire, *Client) {
	t.Helper()
	c, err := nativeLoad(nativeEnv())
	w := &nativeWire{C: c}
	if err != nil {
		w.Blocked = true
		w.Reason = err.Error()
		w.C.Base = ""
	}
	t.Cleanup(func() {
		status := "pass"
		if t.Failed() {
			status = "no_pass"
		}
		if w.Blocked {
			status = "block"
		}
		dir := os.Getenv("ENGAGELAB_RESULTS_DIR")
		if dir == "" {
			dir = ".test-results"
		}
		if err := os.MkdirAll(dir, 0700); err != nil {
			t.Error("cannot create results directory")
			return
		}
		b, _ := json.MarshalIndent(map[string]interface{}{"run_id": c.Run, "scenario_id": scenario, "language": "Go", "suite": "integration", "environment": c.Env, "base_url": w.C.Base, "status": status, "checks": w.Checks, "reason": w.Reason}, "", "  ")
		if err := os.WriteFile(filepath.Join(dir, scenario+".json"), b, 0600); err != nil {
			t.Error("cannot write results")
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if write && c.Mode != "isolated_write" {
		w.Blocked = true
		w.Reason = "requires isolated_write"
		t.Fatal("block: requires isolated_write")
	}
	return w, w.client()
}
func nativePayload(rid, request string) *PushParam {
	zero := 0
	no := false
	return &PushParam{RequestID: request, To: &PushTo{RegistrationID: []string{rid}}, Body: &PushBody{Platform: []string{"android"}, Notification: &NotificationMessage{Android: &AndroidNotification{Alert: "native SDK test", BadgeSetNum: &zero, IsFold: &no}}, Options: &Options{AutoTruncation: &no}}}
}
func (w *nativeWire) need(t *testing.T, k string) string {
	t.Helper()
	v, e := nativeNeed(nativeEnv(), k)
	if e != nil {
		w.Blocked = true
		w.Reason = e.Error()
		t.Fatal(e)
	}
	return v
}

var nativeCtx = context.Background()
