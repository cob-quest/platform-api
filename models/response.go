package models

type SuccessResponse struct {
	CorId string `json:"corId"` // CorId represents the correlation ID of the attempt.
}
