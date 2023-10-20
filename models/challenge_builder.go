package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChallengeBuilder struct {
	Id            primitive.ObjectID `json:"_id" bson:"_id"`
	CorrelationID string             `json:"cor_id" bson:"cor_id"`
	ChallengeName string             `json:"challengeName" bson:"challenge_name"`
	CreatorName   string             `json:"creatorName" bson:"creator_name"`
	S3Path        string             `json:"s3Path" bson:"s3_path"`
}

type ChallengeBuilderList struct {
	ChallengeBuilders []ChallengeBuilder `json:"challenges" bson:",inline"`
}

