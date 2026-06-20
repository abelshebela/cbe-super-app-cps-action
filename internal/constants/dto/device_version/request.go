package deviceversion

type CreateDeviceVersionRequest struct {
	LatestVersion string `json:"version" bson:"latest_version"`
	Platform      string `json:"platform" bson:"platform"`
	DeviceState   string `json:"device_state" bson:"device_state"`
	ReleaseNotes  string `json:"release_notes" bson:"release_notes"`
}

type UpdateDeviceVersionRequest struct {
	ID            string `json:"_id" bson:"_id"`
	LatestVersion string `json:"version" bson:"latest_version"`
	Platform      string `json:"platform" bson:"platform"`
	DeviceState   string `json:"device_state" bson:"device_state"`
	ReleaseNotes  string `json:"release_notes" bson:"release_notes"`
	Enabled       *bool  `json:"enabled" bson:"enabled"`
}

type SetDeviceStateRequest struct {
	DeviceState string `json:"device_state" bson:"device_state"`
}

type EnableOrDisableDeviceVersion struct {
	ID     string `json:"_id" bson:"_id"`
	Enable bool   `json:"enable" bson:"enable"`
}
