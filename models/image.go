package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
    _Id               primitive.ObjectID `json:"_id" bson:"_id"`
	CorId             string		     `json:"corId" bson:"corId"`
	CreatorName       string             `json:"creatorName" bson:"creatorName"`
	ImageName         string             `json:"imageName" bson:"imageName"`
	ImageTag          string             `json:"imageTag" bson:"imageTag"`
	ImageRegistryLink string             `json:"imageRegistryLink" bson:"imageRegistryLink"`
}