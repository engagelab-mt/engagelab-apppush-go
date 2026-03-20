package engagelab

// DeviceStatusGetParam is the request for querying device online status.
type DeviceStatusGetParam struct {
	RegistrationIDs []string `json:"registration_ids"`
}

// DeviceStatusGetResult is the result for a single device status.
type DeviceStatusGetResult struct {
	RegistrationID string `json:"regid,omitempty"`
	Online         *bool  `json:"online,omitempty"`
	LastOnlineTime string `json:"last_online_time,omitempty"`
}

// DeviceGetResult is the device info (tags + alias).
type DeviceGetResult struct {
	Tags  []string `json:"tags,omitempty"`
	Alias string   `json:"alias,omitempty"`
}

// DeviceSetParam sets tags and alias for a device.
type DeviceSetParam struct {
	Tags  *DeviceSetTags `json:"tags,omitempty"`
	Alias string         `json:"alias,omitempty"`
}

type DeviceSetTags struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// TagsGetResult is the response for listing all tags.
type TagsGetResult struct {
	Tags []string `json:"tags"`
}

// TagSetParam sets registration_ids for a tag.
type TagSetParam struct {
	RegistrationIDs *TagRegistrationIDs `json:"registration_ids"`
}

type TagRegistrationIDs struct {
	Add    []string `json:"add,omitempty"`
	Remove []string `json:"remove,omitempty"`
}

// TagsCountGetResult is the tag count result.
type TagsCountGetResult struct {
	TagsCount map[string]int64 `json:"tagsCount"`
}

// TagQuotaGetResult is the tag/alias quota information.
type TagQuotaGetResult struct {
	Data *TagQuotaData `json:"data,omitempty"`
}

type TagQuotaData struct {
	TotalTagQuota      int64                `json:"totalTagQuota"`
	UseTagQuota        int64                `json:"useTagQuota"`
	TotalAliasQuota    int64                `json:"totalAliasQuota"`
	UseAliasQuota      int64                `json:"useAliasQuota"`
	TagUidQuotaDetails []TagUidQuotaDetail  `json:"tagUidQuotaDetail,omitempty"`
}

type TagUidQuotaDetail struct {
	TagName       string `json:"tagName"`
	TotalUidQuota int64  `json:"totalUidQuota"`
	UseUidQuota   int64  `json:"useUidQuota"`
}

// AliasStatusGetResult is the response for querying alias registration_ids.
type AliasStatusGetResult struct {
	RegistrationIDs []string `json:"registration_ids,omitempty"`
}
