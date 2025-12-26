package customersegmentaion

import (
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func MapToCustomerSegUpdate(c *model.CustomerSegmentation) bson.M {
	update := bson.M{}

	if c.CustomerRole != "" {
		update["customer_role"] = c.CustomerRole
	}
	if c.CustomerSegment != "" {
		update["customer_segment"] = c.CustomerSegment
	}
	if c.CustomerSubSegment != "" {
		update["customer_sub_segment"] = c.CustomerSubSegment
	}
	if c.CustomerGroup != "" {
		update["customer_group"] = c.CustomerGroup
	}

	return update
}
