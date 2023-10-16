package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ImageBuilder struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	CorrelationID string             `json:"cor_id" bson:"cor_id"`
	ImageName     string             `json:"imageName" bson:"image_name"`
	CreatorName   string             `json:"creatorName" bson:"creator_name"`
	S3Path        string             `json:"s3Path" bson:"s3_path"`
}

type ImageBuilderList struct {
	ImageBuilders []ImageBuilder `json:"images" bson:",inline"`
}

