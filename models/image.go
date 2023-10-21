package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	CorrelationID string		 `json:"cor_id" bson:"cor_id"`
	CreatorName     string       `json:"creatorName" bson:"creator_name"`
	ImageName string             `json:"image_name" bson:"image_name"`
	ContainerUrl  string         `json:"containerUrl" bson:"container_url"`
	S3Path		string			`json:"s3Path" bson:"s3_path"`		
}

type ImageList struct {
	Images []Image `json:"images" bson:",inline"`
}
