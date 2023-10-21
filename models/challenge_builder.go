package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChallengeBuilder struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	CorID 		  string             `json:"corId" bson:"corId"`
	ChallengeName string             `json:"challengeName" bson:"challengeName"`
	CreatorName   string             `json:"creatorName" bson:"creatorName"`
	S3Path        string             `json:"s3Path" bson:"s3Path"`
}

type ChallengeBuilderList struct {
	ChallengeBuilders []ChallengeBuilder `json:"challenges" bson:",inline"`
}

