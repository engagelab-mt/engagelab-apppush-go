package engagelab

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"
)

type fullFixture struct {
	ID, Label, Op, Args, Response, Body, Expected, Extra string
	Status                                               int
	Timeout, Delay                                       float64
}
type fullRecord struct {
	Method, Path               string
	Query                      url.Values
	Body                       interface{}
	Raw                        []byte
	Auth, GroupAuth, Multipart bool
}

func fullJSON(raw string) interface{} {
	if raw == "" {
		return nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		panic("invalid compiled fixture")
	}
	return v
}
func fullTree(value interface{}) interface{} {
	b, e := json.Marshal(value)
	if e != nil {
		panic(e)
	}
	return fullJSON(string(b))
}
func fullKey(k string) string {
	k = strings.ToLower(regexp.MustCompile(`[^A-Za-z0-9]`).ReplaceAllString(k, ""))
	switch k {
	case "targets":
		return "target"
	case "clicks":
		return "click"
	}
	return k
}
func fullContains(a, b interface{}) bool {
	switch want := b.(type) {
	case map[string]interface{}:
		actual, ok := a.(map[string]interface{})
		if !ok {
			return false
		}
		keys := map[string]interface{}{}
		for k, v := range actual {
			keys[fullKey(k)] = v
		}
		for k, v := range want {
			got, present := keys[fullKey(k)]
			if !present || !fullContains(got, v) {
				return false
			}
		}
		return true
	case []interface{}:
		actual, ok := a.([]interface{})
		if !ok || len(actual) != len(want) {
			return false
		}
		for i := range want {
			if !fullContains(actual[i], want[i]) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a, b)
	}
}
func fullWriteResult(t *testing.T, id, layer, status, reason string, checks []string) {
	t.Helper()
	dir := os.Getenv("ENGAGELAB_FULL_RESULTS_DIR")
	if dir == "" {
		dir = ".test-results/full"
	}
	dir = filepath.Join(dir, layer)
	if e := os.MkdirAll(dir, 0700); e != nil {
		t.Error("result directory unavailable")
		return
	}
	data, _ := json.MarshalIndent(map[string]interface{}{"scenario_id": id, "language": "Go", "layer": layer, "status": status, "reason": reason, "checks": checks}, "", "  ")
	if e := os.WriteFile(filepath.Join(dir, id+".json"), data, 0600); e != nil {
		t.Error("result write failed")
	}
}

type fullRoundTripper func(*http.Request) (*http.Response, error)

func (fn fullRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return fn(r) }
func fullWav(t *testing.T) string {
	t.Helper()
	b := make([]byte, 1644)
	copy(b, "RIFF")
	binary.LittleEndian.PutUint32(b[4:], 1636)
	copy(b[8:], "WAVEfmt ")
	binary.LittleEndian.PutUint32(b[16:], 16)
	binary.LittleEndian.PutUint16(b[20:], 1)
	binary.LittleEndian.PutUint16(b[22:], 1)
	binary.LittleEndian.PutUint32(b[24:], 8000)
	binary.LittleEndian.PutUint32(b[28:], 16000)
	binary.LittleEndian.PutUint16(b[32:], 2)
	binary.LittleEndian.PutUint16(b[34:], 16)
	copy(b[36:], "data")
	binary.LittleEndian.PutUint32(b[40:], 1600)
	p := filepath.Join(t.TempDir(), "voice.wav")
	if e := os.WriteFile(p, b, 0600); e != nil {
		t.Fatal("cannot create native voice fixture")
	}
	return p
}
func fullContract(t *testing.T, f fullFixture) {
	t.Helper()
	checks := []string{}
	t.Cleanup(func() {
		status := "pass"
		if t.Failed() {
			status = "no_pass"
		}
		fullWriteResult(t, f.ID, "contract", status, "", checks)
	})
	ok := func(v bool, label string) {
		t.Helper()
		if !v {
			t.Fatal(label)
		}
		checks = append(checks, label)
	}
	if f.Extra == `"accessors"` {
		t.Fatal("Java-only accessor fixture cannot run in Go")
	}
	var mu sync.Mutex
	records := []fullRecord{}
	status := f.Status
	response := fullJSON(f.Response)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		var body interface{}
		_ = json.Unmarshal(raw, &body)
		record := fullRecord{r.Method, r.URL.Path, r.URL.Query(), body, raw, r.Header.Get("Authorization") == "Basic "+base64.StdEncoding.EncodeToString([]byte("test-key:test-secret")), r.Header.Get("Authorization") == "Basic "+base64.StdEncoding.EncodeToString([]byte("group-group-key:group-secret")), strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data; boundary=")}
		mu.Lock()
		records = append(records, record)
		resp := response
		code := status
		mu.Unlock()
		if resp == "ECHO" {
			b := body.(map[string]interface{})
			var n int
			fmt.Sscan(b["request_id"].(string), &n)
			resp = map[string]interface{}{"request_id": b["request_id"], "msg_id": fmt.Sprint(n + 100)}
		}
		time.Sleep(time.Duration(f.Delay * float64(time.Second)))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if text, yes := resp.(string); yes {
			_, _ = w.Write([]byte(text))
		} else {
			_ = json.NewEncoder(w).Encode(resp)
		}
	}))
	defer server.Close()
	var args []json.RawMessage
	if e := json.Unmarshal([]byte(f.Args), &args); e != nil {
		t.Fatal("compiled native arguments invalid")
	}
	voice := ""
	if f.Op == "voice_create" {
		voice = fullWav(t)
		args[1], _ = json.Marshal(voice)
	}
	actions := newFullActions(server.URL, "test-key", "test-secret", "group-key", "group-secret", time.Duration(f.Timeout*float64(time.Second)), http.DefaultTransport)
	extra := fullJSON(f.Extra)
	if extra == "host" {
		var dc string
		_ = json.Unmarshal(args[0], &dc)
		wanted := string(Singapore)
		if dc == "japan" {
			wanted = string(Japan)
		}
		if dc == "brazil" {
			wanted = string(Brazil)
		}
		captured := ""
		transport := fullRoundTripper(func(r *http.Request) (*http.Response, error) {
			captured = r.URL.String()
			return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(`{"request_id":"r1","msg_id":"123"}`)), Request: r}, nil
		})
		opts := []Option{WithHTTPClient(&http.Client{Transport: transport})}
		if dc != "default" {
			opts = append(opts, WithDataCenter(DataCenter(wanted)))
		}
		c := NewClient("test-key", "test-secret", opts...)
		_, e := c.Push.Validate(context.Background(), &PushParam{})
		ok(e == nil && captured == wanted+"/v4/push/validate", "configured_datacenter_request")
		return
	}
	if extra == "concurrent" {
		var wg sync.WaitGroup
		results := make([]interface{}, 20)
		errs := make([]error, 20)
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				p := fmt.Sprintf(`{"request_id":"%d","to":{"registration_id":["rid%d"]},"body":{"platform":["android"],"notification":{"alert":"body%d"}}}`, i, i, i)
				results[i], errs[i] = actions.call("validate", []json.RawMessage{json.RawMessage(p)})
			}(i)
		}
		wg.Wait()
		ok(len(records) == 20, "twenty_real_requests")
		for i, r := range results {
			ok(errs[i] == nil && fullContains(fullTree(r), map[string]interface{}{"request_id": fmt.Sprint(i), "msg_id": fmt.Sprint(i + 100)}), "concurrent_response_isolation")
		}
		for _, r := range records {
			b := r.Body.(map[string]interface{})
			id := b["request_id"].(string)
			ok(r.Auth && fullContains(b, map[string]interface{}{"to": map[string]interface{}{"registration_id": []interface{}{"rid" + id}}, "body": map[string]interface{}{"notification": map[string]interface{}{"alert": "body" + id}}}), "concurrent_request_isolation")
		}
		return
	}

	if extra == "logs" {
		capture, e := os.CreateTemp(t.TempDir(), "default-logs-*")
		ok(e == nil, "capture_file_available")
		var successfulErr, failure error
		func() {
			previousOut, previousErr, previousLogger := os.Stdout, os.Stderr, log.Writer()
			defer func() { os.Stdout = previousOut; os.Stderr = previousErr; log.SetOutput(previousLogger) }()
			os.Stdout = capture
			os.Stderr = capture
			log.SetOutput(capture)
			fmt.Fprintln(os.Stdout, "capture-active-stdout")
			log.Print("capture-active-stderr")
			_, successfulErr = actions.call(f.Op, args)
			mu.Lock()
			status = 400
			response = map[string]interface{}{"error": map[string]interface{}{"code": 21003, "message": "fixture-invalid"}}
			mu.Unlock()
			_, failure = actions.call(f.Op, args)
		}()
		capture.Close()
		raw, readErr := os.ReadFile(capture.Name())
		ok(readErr == nil, "captured_streams_readable")
		api, structured := failure.(*ApiError)
		ok(successfulErr == nil && structured && api.StatusCode == 400 && api.ErrorBody.Code == 21003, "success_and_error_controls")
		text := string(raw) + fmt.Sprint(failure)
		ok(strings.Contains(text, "capture-active-stdout") && strings.Contains(text, "capture-active-stderr"), "stdout_stderr_observers_active")
		ok(!strings.Contains(text, "test-secret") && !strings.Contains(text, base64.StdEncoding.EncodeToString([]byte("test-key:test-secret"))), "default_logs_and_exceptions_credentials_redacted")
		ok(len(records) == 2 && records[0].Auth && records[1].Auth && records[0].Method == "POST" && records[1].Method == "POST", "two_authenticated_public_API_requests")
		return
	}
	start := time.Now()
	actual, err := actions.call(f.Op, args)
	if extra == "timeout" {
		ok(err != nil, "timeout_error")
		ok(time.Since(start) < 650*time.Millisecond, "bounded_timeout")
		return
	}
	if f.Status >= 400 {
		api, yes := err.(*ApiError)
		ok(yes, "structured_api_error")
		want := fullJSON(f.Response).(map[string]interface{})["error"].(map[string]interface{})
		ok(api.StatusCode == f.Status, "error_http_status")
		ok(float64(api.ErrorBody.Code) == want["code"] && api.ErrorBody.Message == want["message"], "business_code_and_message")
	} else {
		ok(err == nil, "SDK_success")
	}
	mu.Lock()
	r := records[len(records)-1]
	mu.Unlock()
	ok(len(records) > 0, "HTTP_observed")
	ok((f.Op == "group" && r.GroupAuth) || (f.Op != "group" && r.Auth), "authentication")
	wantPath := fullMethodPaths[f.Op]
	pattern := regexp.MustCompile(`\{[^}]+\}`).ReplaceAllString(strings.TrimRight(wantPath[1], "/"), `[^/]+`)
	matched, _ := regexp.MatchString("^"+pattern+"$", strings.TrimRight(r.Path, "/"))
	ok(r.Method == wantPath[0] && matched, "method_and_path")
	if f.Body != "null" && f.Body != "" {
		ok(fullContains(r.Body, fullJSON(f.Body)), "request_fields_exact")
	}
	if f.Expected != "null" && f.Expected != "" {
		ok(fullContains(fullTree(actual), fullJSON(f.Expected)), "parsed_fields_exact")
	}
	if extra == "no-tags" {
		b := r.Body.(map[string]interface{})
		ok(b["tags"] == nil, "no_tag_clearing_instruction")
	}
	if extra == "multipart" {
		b, _ := os.ReadFile(voice)
		ok(r.Multipart && bytes.Contains(r.Raw, b) && bytes.Contains(r.Raw, []byte("en")), "multipart_file_bytes")
	}
	if parts, yes := extra.([]interface{}); yes && parts[0] == "query" {
		ok(fullContains(fullTree(r.Query), parts[1]), "query_values")
	}

}
