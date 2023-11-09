package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Attempt struct {
	_Id               primitive.ObjectID `json:"_id" bson:"_id"`
	ChallengeName     string             `json:"challengeName" bson:"challengeName"`
	CreatorName       string             `json:"creatorName" bson:"creatorName"`
	Participant       string             `json:"participant" bson:"participant"`
	Token             string             `json:"token" bson:"token"`
	ImageRegistryLink string             `json:"imageRegistryLink" bson:"imageRegistryLink"`
	SSHkey            string             `json:"sshkey" bson:"sshkey"`
	Result            float64             `json:"result" bson:"result"`
	IpAddress         string             `json:"ipaddress" bson:"ipaddress"`
	Port              string             `json:"port" bson:"port"`
}

// POST Handler Body
type AttemptSubmitBody struct {
	Token string `json:"token" validate:"required"`
	Result float64 `json:"result" validate:"required"`
}

// AttemptBody is the request body for an attempt
//	@Description	AttemptBody is used to validate the request body for starting or getting an attempt.
//	@Name			AttemptBody
type AttemptBody struct {
	Token string `json:"token" validate:"required"`
	// Email string `json:"email" validate:"required"`
	ChallengeName     string `json:"challengeName" bson:"challengeName"`
	CreatorName       string `json:"creatorName" bson:"creatorName"`
	Participant       string `json:"participant" bson:"participant"`
	ImageRegistryLink string `json:"imageRegistryLink" bson:"imageRegistryLink"`
	CorId             string `json:"corId"`
	EventStatus       string `json:"eventStatus"`
}