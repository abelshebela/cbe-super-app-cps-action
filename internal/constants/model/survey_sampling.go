package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type SamplingMethod string

const (
	SamplingHighValue        SamplingMethod = "HIGH_VALUE"
	SamplingBranchBased      SamplingMethod = "BRANCH_BASED"
	SamplingTimeBased        SamplingMethod = "TIME_BASED"
	SamplingCustomerSegment  SamplingMethod = "CUSTOMER_SEGMENT"
	SamplingFirstTransaction SamplingMethod = "FIRST_TRANSACTION"
	SamplingGeographic       SamplingMethod = "GEOGRAPHIC"
)

type SamplingPeriod string

const (
	PeriodDaily   SamplingPeriod = "DAILY"
	PeriodWeekly  SamplingPeriod = "WEEKLY"
	PeriodMonthly SamplingPeriod = "MONTHLY"
)

type DayOfWeek string

const (
	Monday    DayOfWeek = "MON"
	Tuesday   DayOfWeek = "TUE"
	Wednesday DayOfWeek = "WED"
	Thursday  DayOfWeek = "THU"
	Friday    DayOfWeek = "FRI"
	Saturday  DayOfWeek = "SAT"
	Sunday    DayOfWeek = "SUN"
)

type ThresholdType string

const (
	ThresholdPercentage ThresholdType = "PERCENTAGE"
	ThresholdAbsolute   ThresholdType = "ABSOLUTE"
)

type HighValueConfig struct {
	ThresholdType ThresholdType `bson:"threshold_type" json:"threshold_type"`
	Value         int           `bson:"value"          json:"value"`
}

type TimeBasedConfig struct {
	StartTime  string   `bson:"start_time"   json:"start_time"`
	EndTime    string   `bson:"end_time"     json:"end_time"`
	DaysOfWeek []string `bson:"days_of_week" json:"days_of_week"`
}

type GeographicConfig struct {
	Regions   []string `bson:"regions,omitempty"   json:"regions,omitempty"`
	Districts []string `bson:"districts,omitempty" json:"districts,omitempty"`
	Branches  []string `bson:"branches,omitempty"  json:"branches,omitempty"`
}

// Only the field matching Method is populated; the rest are nil/empty.
type SamplingConfig struct {
	HighValue        *HighValueConfig  `bson:"high_value,omitempty"        json:"high_value,omitempty"`
	BranchBased      []string          `bson:"branch_based,omitempty"      json:"branch_based,omitempty"`
	TimeBased        *TimeBasedConfig  `bson:"time_based,omitempty"        json:"time_based,omitempty"`
	CustomerSegment  []string          `bson:"customer_segment,omitempty"  json:"customer_segment,omitempty"`
	FirstTransaction bool              `bson:"first_transaction,omitempty" json:"first_transaction,omitempty"`
	Geographic       *GeographicConfig `bson:"geographic,omitempty"        json:"geographic,omitempty"`
}

// SurveySamplingConfig is stored in the survey_sampling_configs collection.
// SurveyURL is the Xebo link; at runtime append "?respondent_id="+sha256(customer_id).
type SurveySamplingConfig struct {
	ID        bson.ObjectID  `bson:"_id"        json:"id"`
	Name      string         `bson:"name"       json:"name"`
	Method    SamplingMethod `bson:"method"     json:"method"`
	Enabled   bool           `bson:"enabled"    json:"enabled"`
	Config    SamplingConfig `bson:"config"     json:"config"`
	SurveyURL string         `bson:"survey_url" json:"survey_url"`
	CreatedAt time.Time      `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time      `bson:"updated_at" json:"updated_at"`
}
