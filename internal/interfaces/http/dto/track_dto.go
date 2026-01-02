package dto

import "errors"

// CreateTrackRequest represents the HTTP request to create/update a track.
type CreateTrackRequest struct {
	Md5         string `json:"Md5" validate:"required"`
	StreamTitle string `json:"StreamTitle" validate:"required"`
	StreamUrl   string `json:"StreamUrl"`
}

// Validate validates the CreateTrackRequest.
func (r *CreateTrackRequest) Validate() error {
	if r.Md5 == "" {
		return errors.New("md5 is required")
	}
	if len(r.Md5) != 32 {
		return errors.New("md5 must be 32 characters")
	}
	if r.StreamTitle == "" {
		return errors.New("stream_title is required")
	}
	return nil
}

// TrackResponse represents the HTTP response for track operations.
type TrackResponse struct {
	Rotate int `json:"rotate"`
}

// GetTrackResponse represents the HTTP response for getting a track.
type GetTrackResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Rotate    int    `json:"rotate"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	Listeners int    `json:"listeners"`
}
