package mongodb

import (
	"context"

	models "github.com/mongodb-developer/docker-golang-example/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetProcessedRecords(dbCollection *mongo.Collection) (processedRecords []models.ProcessedRecord, err error) {
	cursor, err := dbCollection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return nil, err
	}
	var records []models.MongoRecord
	var processedRecord models.ProcessedRecord

	if err := cursor.All(context.TODO(), &records); err != nil {
		return nil, err
	}

	processedRecords = []models.ProcessedRecord{}

	for _, v := range records {
		processedRecord.CreatedAt = v.CreatedAt
		processedRecord.ObjectId = v.ObjectID.Hex()
		processedRecord.Key = v.Key
		processedRecord.TotalCount = 0
		if len(v.Counts) > 0 {
			for i := 0; i < len(v.Counts); i++ {
				processedRecord.TotalCount = processedRecord.TotalCount + int64(v.Counts[i])
			}
		}
		processedRecords = append(processedRecords, processedRecord)
	}

	return processedRecords, nil
}
