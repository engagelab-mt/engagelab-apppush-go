package engagelab

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// fullActions binds each operation to a public SDK method, not another runtime.
type fullActions struct {
	client *Client
	group  *GroupPushClient
}

func newFullActions(base, key, secret, gkey, gsecret string, timeout time.Duration, transport http.RoundTripper) *fullActions {
	h := &http.Client{Timeout: timeout, Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	opts := []Option{WithBaseURL(base), WithHTTPClient(h)}
	return &fullActions{NewClient(key, secret, opts...), NewGroupPushClient(gkey, gsecret, opts...)}
}
func fullDecode[T any](raw json.RawMessage) (T, error) {
	var v T
	e := json.Unmarshal(raw, &v)
	return v, e
}
func (a *fullActions) call(op string, raw []json.RawMessage) (interface{}, error) {
	ctx := context.Background()
	switch op {
	case "push":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*PushParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Push.Send(ctx, v0)
	case "validate":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*PushParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Push.Validate(ctx, v0)
	case "push_raw":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[interface{}](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Push.SendRaw(ctx, v0)
	case "batch_regid":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*BatchPushParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Push.BatchByRegID(ctx, v0)
	case "batch_alias":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*BatchPushParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Push.BatchByAlias(ctx, v0)
	case "withdraw":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Push.Withdraw(ctx, v0)
	case "device_get":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Device.Get(ctx, v0)
	case "device_set":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[*DeviceSetParam](raw[1])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Device.Set(ctx, v0, v1)
	case "device_delete":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Device.Delete(ctx, v0)
	case "device_status":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*DeviceStatusGetParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Device.GetStatus(ctx, v0)
	case "token":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*DeviceTokenRegisterParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Device.RegisterToken(ctx, v0)
	case "tags":
		if len(raw) != 0 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		return a.client.Tag.List(ctx)
	case "tag_set":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[*TagSetParam](raw[1])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Tag.Set(ctx, v0, v1)
	case "tag_delete":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[[]string](raw[1])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Tag.Delete(ctx, v0, v1)
	case "tag_status":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Tag.GetDeviceStatus(ctx, v0, v1)
	case "tag_count":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[[]string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Tag.GetCount(ctx, v0, v1)
	case "tag_quota":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[[]string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Tag.GetQuota(ctx, v0, v1)
	case "alias_get":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[[]string](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Alias.Get(ctx, v0, v1)
	case "alias_delete":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[[]string](raw[1])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Alias.Delete(ctx, v0, v1)
	case "schedule_create":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*SchedulePushParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Schedule.Create(ctx, v0)
	case "schedule_update":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[*SchedulePushParam](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Schedule.Update(ctx, v0, v1)
	case "schedule_get":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Schedule.Get(ctx, v0)
	case "schedule_delete":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Schedule.Delete(ctx, v0)
	case "schedule_list":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[int](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Schedule.List(ctx, v0)
	case "schedule_msgs":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Schedule.GetMsgIDs(ctx, v0)
	case "status":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[[]string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Status.MessageDetail(ctx, v0)
	case "lifecycle":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[[]string](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Status.MessageLifecycle(ctx, v0, v1)
	case "batch_lifecycle":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[[]string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Status.BatchMessageDetail(ctx, v0)
	case "users":
		if len(raw) != 3 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		v2, err := fullDecode[int](raw[2])
		if err != nil {
			return nil, err
		}
		return a.client.Status.Users(ctx, v0, v1, v2)
	case "plan_status":
		if len(raw) != 3 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[[]string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		v2, err := fullDecode[string](raw[2])
		if err != nil {
			return nil, err
		}
		return a.client.Status.PlanDetail(ctx, v0, v1, v2)
	case "plan_create":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*PushPlanParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Plan.CreateOrUpdate(ctx, v0)
	case "plan_list":
		if len(raw) != 4 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[int](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[int](raw[1])
		if err != nil {
			return nil, err
		}
		v2, err := fullDecode[*int](raw[2])
		if err != nil {
			return nil, err
		}
		v3, err := fullDecode[string](raw[3])
		if err != nil {
			return nil, err
		}
		return a.client.Plan.List(ctx, v0, v1, v2, v3)
	case "plan_msgs":
		if len(raw) != 3 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		v2, err := fullDecode[string](raw[2])
		if err != nil {
			return nil, err
		}
		return a.client.Plan.QueryMsg(ctx, v0, v1, v2)
	case "plan_delete":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Plan.Delete(ctx, v0)
	case "plan_batch_delete":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Plan.BatchDelete(ctx, v0)
	case "voice_create":
		if len(raw) != 2 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		v1, err := fullDecode[string](raw[1])
		if err != nil {
			return nil, err
		}
		return a.client.Voice.Create(ctx, v0, v1)
	case "voices":
		if len(raw) != 0 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		return a.client.Voice.List(ctx)
	case "voice_get":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Voice.Get(ctx, v0)
	case "voice_delete":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[string](raw[0])
		if err != nil {
			return nil, err
		}
		return nil, a.client.Voice.Delete(ctx, v0)
	case "image":
		if len(raw) != 1 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		v0, err := fullDecode[*OppoImageParam](raw[0])
		if err != nil {
			return nil, err
		}
		return a.client.Image.UploadOppo(ctx, v0)
	case "vip":
		if len(raw) != 0 {
			return nil, fmt.Errorf("wrong native argument count")
		}
		return a.client.App.GetVIPStatus(ctx)
	case "group":
		v, e := fullDecode[*PushParam](raw[0])
		if e != nil {
			return nil, e
		}
		r, e := a.group.Send(ctx, v)
		if e != nil {
			return nil, e
		}
		return map[string]interface{}{"group_msgid": r.GroupMsgID, "successes": r.Successes, "errors": r.Errors}, nil
	default:
		return nil, fmt.Errorf("unsupported native operation")
	}
}
