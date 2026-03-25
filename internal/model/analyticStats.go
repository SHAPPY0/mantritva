package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type AnalyticStat struct {
	ID              primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	StatType        string             `json:"stat_type" bson:"stat_type"`
	Year            int                `json:"year" bson:"year"`
	Month           int                `json:"month" bson:"month"`
	Scope           string             `json:"scope" bson:"scope"`
	Data            []map[string]any   `json:"data" bson:"data"`
	MetaData        map[string]any     `json:"metadata" bson:"metadata"`
	AuditFieldsBSON `bson:",inline"`
}
