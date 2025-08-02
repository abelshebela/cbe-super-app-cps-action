package updatedbulkservice

type SubAccessList struct {
	Key            string `json:"key" bson:"key"`
	Enabled        bool   `json:"enabled" bson:"enabled"`
	AccessListName string `json:"accessListName" bson:"accessListName"`
}

type APPAccessList struct {
	Key            string          `json:"key" bson:"key"`
	Enabled        bool            `json:"enabled" bson:"enabled"`
	AccessListName string          `json:"accessListName" bson:"accessListName"`
	SubAccessList  []SubAccessList `json:"subAccessList" bson:"subAccessList"`
	USSDEnabled    bool            `json:"USSDEnabled" bson:"USSDEnabled"`
}
