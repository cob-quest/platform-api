package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Process struct {
	_Id             primitive.ObjectID `json:"_id" bson:"_id"`
	CorId           string             `json:"corId" bson:"corId"`
	Timestamp       string             `json:"timestamp" bson:"timestamp"`
	CreatorName     string             `json:"creatorName,omitempty" bson:"creatorName,omitempty"`
	ChallengeName   string             `json:"challengeName,omitempty" bson:"challengeName,omitempty"`
	ImageName       string             `json:"imageName,omitempty" bson:"imageName,omitempty"`
	ImageTag        string             `json:"imageTag,omitempty" bson:"imageTag,omitempty"`
	Participant     string             `json:"participant,omitempty" bson:"participant,omitempty"`
	Participants    []string           `json:"participants" bson:"participants"`
	Event           string             `json:"event" bson:"event"`
	EventSuccess    string             `json:"eventSuccess" bson:"eventSuccess"`
}