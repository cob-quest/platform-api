package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Process struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	CorrelationID string             `json:"cor_id" bson:"cor_id"`
	EventSuccess  bool               `json:"eventSuccess" bson:"event_success"`
	ImageName     string             `json:"imageName" bson:"image_name"`
    CreatorName   string             `json:"creatorName" bson:"creator_name"`
    S3path        string             `json:"s3Path" bson:"s3_path"`
    Timestamp     primitive.DateTime `json:"timeStamp" bson:"timestamp"`
}

type ProcessList struct {
	Processes []Process `json:"processes" bson:",inline"`
}
