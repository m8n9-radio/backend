package track

// UpsertTrackCommand represents the command to create or update a track.
type UpsertTrackCommand struct {
	ID    string
	Title string
	Cover string
}

// UpsertTrackResult represents the result of upserting a track.
type UpsertTrackResult struct {
	Rotate int
}

// GetTrackQuery represents the query to get a track.
type GetTrackQuery struct {
	ID string
}

// TrackDTO represents a track for external use.
type TrackDTO struct {
	ID        string
	Title     string
	Cover     string
	Rotate    int
	Likes     int
	Dislikes  int
	Listeners int
}
