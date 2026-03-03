package core

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func UpdateCpsUsers(ctx context.Context, client *mongo.Client, oldJobTitle, newJobTitle string) error {
	// Define the MongoDB collection name
	const collectionName = "cps_users"

	// Access the MongoDB collection (replace 'db' with your MongoDB client instance)
	collection := client.Database(collectionName).Collection(collectionName)

	// Define the filter to match documents with the old job title
	filter := bson.M{"job_title": oldJobTitle}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return err
	}

	// Define the update operation to set the new job title
	update := bson.M{"$set": bson.M{"job_title": newJobTitle}}
	fmt.Println(filter, update, results, "///////////")
	// Perform the update operation
	result, err := collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	// Log the number of documents updated (optional)
	fmt.Printf("Updated %d documents in collection %s\n", result.ModifiedCount, collectionName)

	return nil
}
