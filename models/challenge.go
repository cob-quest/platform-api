package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Challenge struct {
	_Id                 primitive.ObjectID `json:"_id" bson:"_id"`
	CorID               string             `json:"corId" bson:"corId"`
    ChallengeName       string             `json:"challengeName" bson:"challengeName"`
	CreatorName         string             `json:"creatorName" bson:"creatorName"`
	ImageName           string             `json:"imageName" bson:"imageName"`
	ImageTag            string             `json:"imageTag" bson:"imageTag"`
	ImageRegistryLink   string             `json:"imageRegistryLink" bson:"imageRegistryLink"`
	Duration            int                `json:"duration" bson:"duration"`
	Participants        []string           `json:"participants" bson:"participants"`
}