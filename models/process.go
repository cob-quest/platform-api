package models

// "google.golang.org/genproto/googleapis/type/datetime"

type Process struct {
	Timestamp     map[string]interface{} `json:"timestamp" bson:"timestamp"`
	CorId         *string                `json:"corId" bson:"corId"`
	Event         *string                `json:"event" bson:"event"`
	EventStatus   *string                `json:"eventStatus" bson:"eventStatus"`
	CreatorName   *string                `json:"creatorName,omitempty" bson:"creatorName,omitempty"`
	ChallengeName *string                `json:"challengeName,omitempty" bson:"challengeName,omitempty"`
	ImageName     *string                `json:"imageName,omitempty" bson:"imageName,omitempty"`
	ImageTag      *string                `json:"imageTag,omitempty" bson:"imageTag,omitempty"`
	Participant   *string                `json:"participant,omitempty" bson:"participant,omitempty"`
	Participants  *[]string              `json:"participants,omitempty" bson:"participants,omitempty"`
}