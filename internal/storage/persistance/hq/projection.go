package hq

import "go.mongodb.org/mongo-driver/v2/bson"

func HQProjection() bson.M {
	return bson.M{
		"_id":                    1,
		"latest_ios_version":     1,
		"latest_android_version": 1,
	}
}
