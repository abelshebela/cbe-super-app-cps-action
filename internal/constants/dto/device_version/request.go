package deviceversion

type CreateDeviceVersionRequest struct {
	LatestVersion string `json:"version" bson:"latest_version"`
	Platform      string `json:"platform" bson:"platform"`
	ForceUpdate   bool   `json:"force_update" bson:"force_update"`
	ReleaseNotes  string `json:"release_notes" bson:"release_notes"`
}

type UpdateDeviceVersionRequest struct {
	ID            string `json:"_id" bson:"_id"`
	LatestVersion string `json:"version" bson:"latest_version"`
	Platform      string `json:"platform" bson:"platform"`
	ForceUpdate   bool   `json:"force_update" bson:"force_update"`
	ReleaseNotes  string `json:"release_notes" bson:"release_notes"`
}

type EnableOrDisableDeviceVersion struct {
	ID     string `json:"_id" bson:"_id"`
	Enable bool   `json:"enable" bson:"enable"`
}
