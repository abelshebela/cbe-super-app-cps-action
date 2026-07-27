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
	SamplingFrequency        SamplingMethod = "FREQUENCY_CONTROLLED"
	SamplingPercentageBased  SamplingMethod = "PERCENTAGE_BASED"
	SamplingGeographic       SamplingMethod = "GEOGRAPHIC"
)

type SamplingPeriod string

const (
	PeriodDaily   SamplingPeriod = "daily"
	PeriodWeekly  SamplingPeriod = "weekly"
	PeriodMonthly SamplingPeriod = "monthly"
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

type HighValueConfig struct {
	MinAmount float64 `bson:"min_amount" json:"min_amount"`
}

type BranchBasedConfig struct {
	BranchCodes []string `bson:"branch_codes" json:"branch_codes"`
}

type TimeBasedConfig struct {
	StartTime  string   `bson:"start_time"   json:"start_time"`   // "HH:MM"
	EndTime    string   `bson:"end_time"     json:"end_time"`     // "HH:MM"
	DaysOfWeek []string `bson:"days_of_week" json:"days_of_week"` // ["MON","TUE",...]
}

type CustomerSegmentConfig struct {
	Segments []string `bson:"segments" json:"segments"` // ["CORPORATE","RETAIL","IFB"]
}

type FirstTransactionConfig struct {
	Period string `bson:"period" json:"period"` // "daily" | "weekly" | "monthly"
}

type FrequencyConfig struct {
	EveryNTransactions int `bson:"every_n_transactions,omitempty" json:"every_n_transactions,omitempty"` // e.g. every 100th txn
	CooldownHours      int `bson:"cooldown_hours,omitempty"      json:"cooldown_hours,omitempty"`        // e.g. 168 = once per week
}

type PercentageBasedConfig struct {
	Percentage float64 `bson:"percentage" json:"percentage"` // e.g. 10.5 = 10.5% of transactions
}

type GeographicConfig struct {
	Regions   []string `bson:"regions,omitempty"   json:"regions,omitempty"`
	Districts []string `bson:"districts,omitempty" json:"districts,omitempty"`
	Branches  []string `bson:"branches,omitempty"    json:"branches,omitempty"`
}

type SamplingConfig struct {
	HighValue        *HighValueConfig        `bson:"high_value,omitempty"        json:"high_value,omitempty"`
	BranchBased      *BranchBasedConfig      `bson:"branch_based,omitempty"      json:"branch_based,omitempty"`
	TimeBased        *TimeBasedConfig        `bson:"time_based,omitempty"        json:"time_based,omitempty"`
	CustomerSegment  string                  `bson:"customer_segment,omitempty"  json:"customer_segment,omitempty"`
	FirstTransaction *FirstTransactionConfig `bson:"first_transaction,omitempty" json:"first_transaction,omitempty"`
	Frequency        *FrequencyConfig        `bson:"frequency,omitempty"         json:"frequency,omitempty"`
	PercentageBased  *PercentageBasedConfig  `bson:"percentage_based,omitempty"  json:"percentage_based,omitempty"`
	Geographic       *GeographicConfig       `bson:"geographic,omitempty"        json:"geographic,omitempty"`
}

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
