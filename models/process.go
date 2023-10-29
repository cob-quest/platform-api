package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Process struct {
	_Id          primitive.ObjectID `json:"_id" bson:"_id"`
	CorId        string             `json:"corId" bson:"corId"`
	CreatorName  string             `json:"creatorName" bson:"creatorName"`
	Event        string             `json:"event" bson:"event"`
	EventSuccess bool               `json:"eventSuccess" bson:"eventSuccess"`
	Timestamp    time.Time          `json:"timestamp" bson:"timestamp"`
	ImageName    string             `json:"imageName" bson:"imageName"`
	S3Path       string             `json:"s3Path,omitempty" bson:"s3Path,omitempty"`
	Duration     int                `json:"duration,omitempty" bson:"duration,omitempty"`
	Participants []string           `json:"participants,omitempty" bson:"participants,omitempty"`
	AssignmentId string             `json:"assignmentId,omitempty" bson:"assignmentId,omitempty"`
}

type ProcessList struct {
	Processes []Process `json:"processes" bson:",inline"`
}
