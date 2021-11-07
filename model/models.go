package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MongoRecord struct {
	ObjectID  primitive.ObjectID `bson:"_id" json:"_id"`
	Key       string             `json:"key"`
	CreatedAt time.Time          `json:"createdAt"`
	Counts    []int32            `json:"counts"`
}

type ProcessedRecord struct {
	ObjectId   string
	Key        string
	CreatedAt  time.Time
	TotalCount int64
}

type RelatedRecordInput struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type RelatedRecordsResponse struct {
	Code         int             `json:"code"`
	Message      string          `json:"message"`
	Records      []RelatedRecord `json:"records"`
	ErrorDetails string          `json:"errorDetails,omitempty"`
}

type RelatedRecord struct {
	Key        string `json:"code"`
	CreatedAt  string `json:"message"`
	TotalCount int    `json:"records"`
}

type InMemoryResponse struct {
	Message string `json:"message,omitempty"`
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
}

type InMemoryRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
