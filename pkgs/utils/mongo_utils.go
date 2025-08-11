package utils

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MongoIDConverter provides utilities for MongoDB ID handling
type MongoIDConverter struct{}

// ConvertStringToObjectID safely converts a string to ObjectID
func (m *MongoIDConverter) ConvertStringToObjectID(id string) (bson.ObjectID, error) {
	return ParsePrimitiveObjectID(id)
}

// ConvertFilterValues converts string values in filters to ObjectIDs where needed
func (m *MongoIDConverter) ConvertFilterValues(filters map[string]interface{}, idFields []string) (map[string]interface{}, error) {
	converted := make(map[string]interface{})

	for key, value := range filters {
		if sliceContains(idFields, key) {
			if strValue, ok := value.(string); ok {
				objID, err := m.ConvertStringToObjectID(strValue)
				if err != nil {
					return nil, fmt.Errorf("invalid %s format: %v", key, err)
				}
				converted[key] = objID
			} else {
				converted[key] = value
			}
		} else {
			converted[key] = value
		}
	}

	return converted, nil
}

// BuildLookupPipeline creates an aggregation pipeline with lookups
func (m *MongoIDConverter) BuildLookupPipeline(matchFilter bson.M, lookups []bson.M, skip, limit int64) []bson.M {
	pipeline := []bson.M{
		{"$match": matchFilter},
		{"$skip": skip},
		{"$limit": limit},
	}

	// Add lookups
	for _, lookup := range lookups {
		pipeline = append(pipeline, lookup)
	}

	return pipeline
}

// Helper function to check if slice contains string
func sliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
