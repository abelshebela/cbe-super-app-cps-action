package updatedbulkservice

type BulkServiceDTO struct {
	Keys []string `json:"keys" bson:"keys"`
}
