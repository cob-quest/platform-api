package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChallengeBuilder struct {
    _Id           primitive.ObjectID `json:"_id" bson:"_id"`
	CorID 		  string             `json:"corId" bson:"corId"`
	ImageName     string             `json:"imageName" bson:"imageName"`
	CreatorName   string             `json:"creatorName" bson:"creatorName"`
    Duration      int                `json:"duration" bson:"duration"`
    Participants  []string           `json:"participants" bson:"participants"`
    AssignmentId  string             `json:"assignmentId" bson:"assignmentId"`
}

type ChallengeBuilderList struct {
	ChallengeBuilders []ChallengeBuilder `json:"challenges" bson:",inline"`
}

