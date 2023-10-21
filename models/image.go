package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	ImageName string             `json:"image_name" bson:"image_name" validate:"required"`
	ImageVer  string             `json:"image_ver" bson:"image_ver" validate:"required"`
}

type ImageList struct {
	Images []Image `json:"images" bson:",inline"`
}
