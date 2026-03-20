package engagelab

// SchedulePushParam is the request for creating a scheduled push.
type SchedulePushParam struct {
	Name    string          `json:"name,omitempty"`
	Enabled *bool           `json:"enabled,omitempty"`
	Trigger *ScheduleTrigger `json:"trigger,omitempty"`
	Push    *PushParam      `json:"push,omitempty"`
}

type ScheduleTrigger struct {
	Single     *TriggerSingle     `json:"single,omitempty"`
	Periodical *TriggerPeriodical `json:"periodical,omitempty"`
}

type TriggerSingle struct {
	Time     string `json:"time,omitempty"`     // "yyyy-MM-dd HH:mm:ss"
	ZoneType *int   `json:"zone_type,omitempty"`
}

type TriggerPeriodical struct {
	Start     string   `json:"start,omitempty"`     // "yyyy-MM-dd HH:mm:ss"
	End       string   `json:"end,omitempty"`       // "yyyy-MM-dd HH:mm:ss"
	Time      string   `json:"time,omitempty"`      // "HH:mm:ss"
	Frequency *int     `json:"frequency,omitempty"`
	TimeUnit  string   `json:"time_unit,omitempty"` // "day", "WEEK", "MONTH"
	Point     []string `json:"point,omitempty"`
	ZoneType  *int     `json:"zone_type,omitempty"`
}

// SchedulePushResult is the response after creating a schedule.
type SchedulePushResult struct {
	ScheduleID string `json:"schedule_id,omitempty"`
	Name       string `json:"name,omitempty"`
}

// SchedulePushGetResult is the response for getting a schedule by ID.
type SchedulePushGetResult struct {
	ScheduleID string          `json:"schedule_id,omitempty"`
	Name       string          `json:"name,omitempty"`
	Enabled    *bool           `json:"enabled,omitempty"`
	Trigger    *ScheduleTrigger `json:"trigger,omitempty"`
	Push       *PushParam      `json:"push,omitempty"`
}

// SchedulePushListResult is the paginated schedule list response.
type SchedulePushListResult struct {
	TotalCount  int                      `json:"total_count"`
	TotalPages  int                      `json:"total_pages"`
	CurrentPage int                      `json:"page"`
	Schedules   []SchedulePushListDetail `json:"schedules,omitempty"`
}

type SchedulePushListDetail struct {
	ScheduleID string      `json:"schedule_id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Enabled    *bool       `json:"enabled,omitempty"`
	Trigger    interface{} `json:"trigger,omitempty"`
	Push       interface{} `json:"push,omitempty"`
}

// SchedulePushDetailGetResult contains msg IDs for a schedule.
type SchedulePushDetailGetResult struct {
	Count   int64         `json:"count"`
	MsgIDs  []interface{} `json:"MsgIds,omitempty"`
}
