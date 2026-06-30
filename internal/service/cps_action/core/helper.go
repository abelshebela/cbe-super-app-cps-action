package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

var VolatileChecksumKeys = map[string]bool{
	"id":          true,
	"action_code": true,
	"version":     true,
	"user_code":   true,

	"is_deleted": true,
	"is_enabled": true,

	"created_at":                    true,
	"create_at":                     true,
	"updated_at":                    true,
	"update_at":                     true,
	"last_modified_at":              true,
	"last_modified":                 true,
	"maker_action_time":             true,
	"approved_at":                   true,
	"reversed_at":                   true,
	"deleted_at":                    true,
	"date_joined":                   true,
	"issued_date":                   true,
	"sent_at":                       true,
	"published_at":                  true,
	"verified_at":                   true,
	"completed_at":                  true,
	"claimed_at":                    true,
	"linked_at":                     true,
	"initial_linked_at":             true,
	"initiated_linked_at":           true,
	"pin_changed_at":                true,
	"password_changed_at":           true,
	"application_installation_date": true,
	"last_login":                    true,
	"last_login_attempt":            true,
	"last_online_date":              true,
	"otp_last_tried_at":             true,
	"otp_last_verified_at":          true,
	"created_at_password_expiry":    true,
	"updated_at_password_expiry":    true,
	"created_at_block":              true,
	"updated_at_block":              true,
	"created_at_archive":            true,
	"updated_at_archive":            true,
	"created_at_total_cap":          true,
	"updated_at_total_cap":          true,
	"next_attempt_count":            true,

	// --- file / image / media URLs (Minio key varies per upload) ---
	"logo":            true,
	"icon":            true,
	"app_icon":        true,
	"company_logo":    true,
	"donation_icon":   true,
	"image":           true,
	"image_url":       true,
	"cover_image":     true,
	"cover_image_url": true,
	"banner_image":    true,
	"photo":           true,
	"selfie_photo":    true,
	"picture":         true,
	"thumbnail":       true,
	"avatar":          true,
	"document_front":  true,
	"document_back":   true,
	"signature":       true,
	"video_url":       true,
	"receipt_link":    true,
	"url":             true,
}

func ComputeActionChecksum(requestAction, uniqueID string, currentAction interface{}) string {
	raw, _ := json.Marshal(currentAction)

	var m map[string]interface{}
	if json.Unmarshal(raw, &m) == nil {
		for k := range VolatileChecksumKeys {
			delete(m, k)
		}
		raw, _ = json.Marshal(m)
	}

	h := sha256.New()
	h.Write([]byte(requestAction))
	h.Write([]byte("|"))
	h.Write([]byte(uniqueID))
	h.Write([]byte("|"))
	h.Write(raw)
	return hex.EncodeToString(h.Sum(nil))
}
