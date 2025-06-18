package utils

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OtpDAL interface {
	DeleteMany(ctx context.Context, filter interface{}) (int64, error)
}

func HandleExpiredOTP(ctx context.Context, otpDal OtpDAL) {
	now := time.Now().UTC()

	expiredQuery := bson.M{
		"expires_at": bson.M{
			"$lt": now,
		},
	}

	log.Println("==== ", expiredQuery)

	deletedCount, err := otpDal.DeleteMany(ctx, expiredQuery)
	if err != nil {
		log.Println("------ error in removing expired otps")
		return
	}

	log.Printf("------ removed expired otps: %d\n", deletedCount)
}

type OtpCollection struct {
	Collection *mongo.Collection
}

func (o *OtpCollection) DeleteMany(ctx context.Context, filter interface{}) (int64, error) {
	res, err := o.Collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
