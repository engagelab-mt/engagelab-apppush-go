// Package engagelab provides a Go client for the EngageLab AppPush REST API.
//
// Create a client with NewClient, then access services via the client fields:
//
//	client := engagelab.NewClient("appKey", "masterSecret")
//
//	// Push a notification
//	result, err := client.Push.Send(ctx, &engagelab.PushParam{...})
//
//	// Query device info
//	device, err := client.Device.Get(ctx, registrationID)
//
// For Group Push, use a separate client:
//
//	gc := engagelab.NewGroupPushClient("groupKey", "groupMasterSecret")
//	result, err := gc.Send(ctx, &engagelab.GroupPushParam{...})
//
// All API errors are returned as *ApiError, which includes the HTTP status code
// and the business error code from EngageLab.
package engagelab
