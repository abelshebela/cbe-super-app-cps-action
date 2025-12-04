package hq

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func BuildHQUpdateDoc(updates map[string]interface{}, now time.Time) bson.M {
	updateDoc := bson.M{}

	for field, value := range updates {
		if value == "" {
			continue
		}
		updateDoc[field] = value
		switch field {
		case "block_time":
			updateDoc["updated_at_block"] = now
		case "archive_time":
			updateDoc["updated_at_archive"] = now
		case "password_expiry":
			updateDoc["updated_at_password_expiry"] = now
		}
	}

	return updateDoc
}
