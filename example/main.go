package main

import (
	"context"
	"fmt"
	"log"

	engagelab "github.com/engagelab-mt/engagelab-apppush-go"
)

func main() {
	client := engagelab.NewClient(
		"your-app-key",
		"your-master-secret",
		engagelab.WithDataCenter(engagelab.Singapore),
	)

	ctx := context.Background()

	pushExample(ctx, client)
	deviceExample(ctx, client)
	tagExample(ctx, client)
	scheduleExample(ctx, client)
	statusExample(ctx, client)
	planExample(ctx, client)
	groupPushExample(ctx)
}

func pushExample(ctx context.Context, client *engagelab.Client) {
	fmt.Println("=== Push Example ===")

	result, err := client.Push.Send(ctx, &engagelab.PushParam{
		From: "push",
		To:   "all",
		Body: &engagelab.PushBody{
			Platform: "all",
			Notification: &engagelab.NotificationMessage{
				Alert: "Hello from Go SDK!",
				Android: &engagelab.AndroidNotification{
					Alert: "Hello Android!",
					Title: "Test Push",
				},
				IOS: &engagelab.IOSNotification{
					Alert: "Hello iOS!",
				},
			},
		},
	})
	if err != nil {
		log.Printf("Push failed: %v", err)
		return
	}
	fmt.Printf("Push sent: msg_id=%s\n", result.MsgID)

	// Push with custom message
	result, err = client.Push.Send(ctx, &engagelab.PushParam{
		To: "all",
		Body: &engagelab.PushBody{
			Platform: "all",
			Message: &engagelab.CustomMessage{
				Title:      "Custom Message",
				MsgContent: "This is a passthrough message",
				Extras: map[string]interface{}{
					"key1": "value1",
				},
			},
		},
	})
	if err != nil {
		log.Printf("Custom push failed: %v", err)
		return
	}
	fmt.Printf("Custom push sent: msg_id=%s\n", result.MsgID)

	// Push to specific targets
	result, err = client.Push.Send(ctx, &engagelab.PushParam{
		To: &engagelab.PushTo{
			Tag:   []string{"vip", "premium"},
			Alias: []string{"user_001"},
		},
		Body: &engagelab.PushBody{
			Platform: []string{"android", "ios"},
			Notification: &engagelab.NotificationMessage{
				Alert: "Targeted push!",
			},
		},
	})
	if err != nil {
		log.Printf("Targeted push failed: %v", err)
		return
	}
	fmt.Printf("Targeted push sent: msg_id=%s\n", result.MsgID)

	// Batch push by registration ID
	batchResult, err := client.Push.BatchByRegID(ctx, &engagelab.BatchPushParam{
		Requests: []engagelab.BatchPushRequest{
			{
				Target:   "regid_001",
				Platform: "android",
				Notification: &engagelab.NotificationMessage{
					Alert: "Hello regid_001!",
				},
			},
			{
				Target:   "regid_002",
				Platform: "android",
				Notification: &engagelab.NotificationMessage{
					Alert: "Hello regid_002!",
				},
			},
		},
	})
	if err != nil {
		log.Printf("Batch push failed: %v", err)
		return
	}
	fmt.Printf("Batch push result: %+v\n", batchResult)

	// Validate push
	_, err = client.Push.Validate(ctx, &engagelab.PushParam{
		To: "all",
		Body: &engagelab.PushBody{
			Platform: "all",
			Notification: &engagelab.NotificationMessage{
				Alert: "Validation test",
			},
		},
	})
	if err != nil {
		log.Printf("Push validate failed: %v", err)
	} else {
		fmt.Println("Push validation passed")
	}
}

func deviceExample(ctx context.Context, client *engagelab.Client) {
	fmt.Println("\n=== Device Example ===")

	device, err := client.Device.Get(ctx, "registration_id_001")
	if err != nil {
		log.Printf("Get device failed: %v", err)
		return
	}
	fmt.Printf("Device tags: %v, alias: %s\n", device.Tags, device.Alias)

	err = client.Device.Set(ctx, "registration_id_001", &engagelab.DeviceSetParam{
		Tags: &engagelab.DeviceSetTags{
			Add:    []string{"new_tag"},
			Remove: []string{"old_tag"},
		},
		Alias: "new_alias",
	})
	if err != nil {
		log.Printf("Set device failed: %v", err)
	}

	statusList, err := client.Device.GetStatus(ctx, &engagelab.DeviceStatusGetParam{
		RegistrationIDs: []string{"regid_001", "regid_002"},
	})
	if err != nil {
		log.Printf("Get device status failed: %v", err)
		return
	}
	for _, s := range statusList {
		fmt.Printf("Device %s online: %v\n", s.RegistrationID, s.Online)
	}
}

func tagExample(ctx context.Context, client *engagelab.Client) {
	fmt.Println("\n=== Tag Example ===")

	tags, err := client.Tag.List(ctx)
	if err != nil {
		log.Printf("List tags failed: %v", err)
		return
	}
	fmt.Printf("Tags: %v\n", tags.Tags)

	err = client.Tag.Set(ctx, "vip", &engagelab.TagSetParam{
		RegistrationIDs: &engagelab.TagRegistrationIDs{
			Add: []string{"regid_001", "regid_002"},
		},
	})
	if err != nil {
		log.Printf("Set tag failed: %v", err)
	}

	count, err := client.Tag.GetCount(ctx, []string{"vip"}, []string{"android"})
	if err != nil {
		log.Printf("Get tag count failed: %v", err)
		return
	}
	fmt.Printf("Tag counts: %v\n", count.TagsCount)

	alias, err := client.Alias.Get(ctx, "user_001", []string{"android", "ios"})
	if err != nil {
		log.Printf("Get alias failed: %v", err)
		return
	}
	fmt.Printf("Alias registration_ids: %v\n", alias.RegistrationIDs)
}

func scheduleExample(ctx context.Context, client *engagelab.Client) {
	fmt.Println("\n=== Schedule Example ===")

	enabled := true
	schedResult, err := client.Schedule.Create(ctx, &engagelab.SchedulePushParam{
		Name:    "Daily Push",
		Enabled: &enabled,
		Trigger: &engagelab.ScheduleTrigger{
			Single: &engagelab.TriggerSingle{
				Time: "2026-12-31 10:00:00",
			},
		},
		Push: &engagelab.PushParam{
			To: "all",
			Body: &engagelab.PushBody{
				Platform: "all",
				Notification: &engagelab.NotificationMessage{
					Alert: "Scheduled push!",
				},
			},
		},
	})
	if err != nil {
		log.Printf("Create schedule failed: %v", err)
		return
	}
	fmt.Printf("Schedule created: id=%s\n", schedResult.ScheduleID)

	list, err := client.Schedule.List(ctx, 1)
	if err != nil {
		log.Printf("List schedules failed: %v", err)
		return
	}
	fmt.Printf("Total schedules: %d\n", list.TotalCount)
}

func statusExample(ctx context.Context, client *engagelab.Client) {
	fmt.Println("\n=== Status Example ===")

	users, err := client.Status.Users(ctx, "day", "2026-03-01", 7)
	if err != nil {
		log.Printf("Get user status failed: %v", err)
		return
	}
	fmt.Printf("User stats: %d items\n", len(users.Items))

	msgStats, err := client.Status.MessageDetail(ctx, []string{"msg_001", "msg_002"})
	if err != nil {
		log.Printf("Get message status failed: %v", err)
		return
	}
	for id, stat := range msgStats {
		fmt.Printf("Message %s: sent=%d delivered=%d\n", id, stat.Sent, stat.Delivered)
	}
}

func planExample(ctx context.Context, client *engagelab.Client) {
	fmt.Println("\n=== Plan Example ===")

	planResult, err := client.Plan.CreateOrUpdate(ctx, &engagelab.PushPlanParam{
		PlanID:          "marketing_plan_001",
		PlanDescription: "Spring promotion campaign",
	})
	if err != nil {
		log.Printf("Create plan failed: %v", err)
		return
	}
	fmt.Printf("Plan created: id=%s\n", planResult.PlanID)

	planList, err := client.Plan.List(ctx, 1, 10, nil, "")
	if err != nil {
		log.Printf("List plans failed: %v", err)
		return
	}
	fmt.Printf("Total plans: %d\n", planList.Total)
}

func groupPushExample(ctx context.Context) {
	fmt.Println("\n=== Group Push Example ===")

	groupClient := engagelab.NewGroupPushClient(
		"your-group-key",
		"your-group-master-secret",
		engagelab.WithDataCenter(engagelab.Singapore),
	)

	result, err := groupClient.Send(ctx, &engagelab.PushParam{
		To: "all",
		Body: &engagelab.PushBody{
			Platform: "all",
			Notification: &engagelab.NotificationMessage{
				Alert: "Group push to all apps!",
			},
		},
	})
	if err != nil {
		log.Printf("Group push failed: %v", err)
		return
	}
	fmt.Printf("Group push sent: group_msgid=%s\n", result.GroupMsgID)
}
