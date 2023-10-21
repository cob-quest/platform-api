package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Process struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	CorId 		  string             `json:"corId" bson:"corId"`
	EventSuccess  bool               `json:"eventSuccess" bson:"eventSuccess"`
	ImageName     string             `json:"imageName" bson:"imageName"`
    CreatorName   string             `json:"creatorName" bson:"creatorName"`
    S3Path        string             `json:"s3Path" bson:"s3Path"`
    Timestamp     primitive.DateTime `json:"timestamp" bson:"timestamp"`
}

type ProcessList struct {
	Processes []Process `json:"processes" bson:",inline"`
}
