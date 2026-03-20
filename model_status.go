package engagelab

// UserStatusGetResult is the response for user statistics.
type UserStatusGetResult struct {
	TimeUnit string           `json:"TimeUnit,omitempty"`
	Start    string           `json:"start,omitempty"`
	Duration int              `json:"duration,omitempty"`
	Items    []UserStatusItem `json:"items,omitempty"`
}

type UserStatusItem struct {
	Time    string              `json:"time,omitempty"`
	Android *UserStatusPlatform `json:"android,omitempty"`
	IOS     *UserStatusPlatform `json:"ios,omitempty"`
}

type UserStatusPlatform struct {
	New    int64 `json:"new"`
	Active int64 `json:"active"`
	Online int64 `json:"online"`
}

// MessageStatusGetResult is the response for message delivery statistics.
type MessageStatusGetResult struct {
	Targets     int64              `json:"targets"`
	Sent        int64              `json:"sent"`
	Delivered   int64              `json:"delivered"`
	Impressions int64              `json:"impressions"`
	Clicks      int64              `json:"clicks"`
	Sub         *MessageStatusSub  `json:"sub,omitempty"`
}

type MessageStatusSub struct {
	Notification *MessageStatusDetail `json:"notification,omitempty"`
	Message      *MessageStatusDetail `json:"message,omitempty"`
}

type MessageStatusDetail struct {
	Target      int64                      `json:"target"`
	Sent        int64                      `json:"sent"`
	Delivered   int64                      `json:"delivered"`
	Impressions int64                      `json:"impressions"`
	Click       int64                      `json:"click"`
	SubAndroid  *MessageStatusAndroid      `json:"sub_android,omitempty"`
	SubIOS      *MessageStatusIOS          `json:"sub_ios,omitempty"`
}

type MessageStatusAndroid struct {
	EngageLabAndroid *MessageStatusChannel `json:"engageLab_android,omitempty"`
	Xiaomi           *MessageStatusChannel `json:"xiaomi,omitempty"`
	Huawei           *MessageStatusChannel `json:"huawei,omitempty"`
	Honor            *MessageStatusChannel `json:"honor,omitempty"`
	Meizu            *MessageStatusChannel `json:"meizu,omitempty"`
	Oppo             *MessageStatusChannel `json:"oppo,omitempty"`
	Vivo             *MessageStatusChannel `json:"vivo,omitempty"`
	FCM              *MessageStatusChannel `json:"fcm,omitempty"`
}

type MessageStatusIOS struct {
	EngageLabIOS *MessageStatusChannel `json:"engageLab_ios,omitempty"`
	APNS         *MessageStatusChannel `json:"apns,omitempty"`
}

type MessageStatusChannel struct {
	Targets     int64 `json:"targets"`
	Sent        int64 `json:"sent"`
	Delivered   int64 `json:"delivered"`
	Impressions int64 `json:"impressions"`
	Clicks      int64 `json:"clicks"`
}

// MessageLifecycleGetResult is the lifecycle status for a message on a device.
type MessageLifecycleGetResult struct {
	Status       string `json:"status,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
