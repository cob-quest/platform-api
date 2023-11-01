package models

import (
	// "google.golang.org/genproto/googleapis/type/datetime"
)

type Process struct {
	Timestamp     map[string]interface{} `json:"timestamp" bson:"timestamp"`
	CorId         *string           `json:"corId" bson:"corId"`
	Event         *string           `json:"event" bson:"event"`
	EventStatus   *string           `json:"eventStatus" bson:"eventStatus"`
	CreatorName   *string           `json:"creatorName" bson:"creatorName,omitempty"`
	ChallengeName *string           `json:"challengeName" bson:"challengeName,omitempty"`
	ImageName     *string           `json:"imageName" bson:"imageName,omitempty"`
	ImageTag      *string           `json:"imageTag" bson:"imageTag,omitempty"`
	Participant   *string           `json:"participant" bson:"participant,omitempty"`
	Participants  *[]string         `json:"participants" bson:"participants,omitempty"`
}