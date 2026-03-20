package engagelab

// PushPlanParam is the request for creating or updating a push plan.
type PushPlanParam struct {
	PlanID          string `json:"plan_id,omitempty"`
	PlanDescription string `json:"plan_description,omitempty"`
}

// PushPlanResult is the response after creating/updating a push plan.
type PushPlanResult struct {
	PlanID string `json:"plan_id,omitempty"`
}

// PushPlanDeleteResult is the response after deleting a push plan.
type PushPlanDeleteResult struct {
	PlanID string `json:"plan_id,omitempty"`
}

// PushPlanListResult is the paginated list of push plans.
type PushPlanListResult struct {
	PushPlanInfo []PushPlanInfo `json:"push_plan_info,omitempty"`
	Total        int            `json:"total,omitempty"`
}

type PushPlanInfo struct {
	PushID          string `json:"push_id,omitempty"`
	PlanDescription string `json:"plan_description,omitempty"`
	Count           int    `json:"count,omitempty"`
	CreateTime      int64  `json:"create_time,omitempty"`
	LastUsedTime    int64  `json:"last_used_time,omitempty"`
}

// PushPlanMsgQueryResult maps plan IDs to their message IDs.
type PushPlanMsgQueryResult map[string]PlanMsgInfo

type PlanMsgInfo struct {
	MsgIDs []string `json:"msg_ids,omitempty"`
}
