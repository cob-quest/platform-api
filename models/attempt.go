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

type AttemptList struct {
	Attempts []Attempt `json:"attempts" bson:",inline"`
}
