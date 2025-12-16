package dto

type CreateTrackRequest struct {
	Md5         string `json:"Md5" validate:"required"`
	StreamTitle string `json:"StreamTitle" validate:"required"`
	StreamUrl   string `json:"StreamUrl"`
}

type TrackResponse struct {
	Rotate int `json:"rotate"`
}
