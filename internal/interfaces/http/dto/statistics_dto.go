package dto

// TrackStats represents track statistics in HTTP response.
type TrackStats struct {
	Title     string `json:"title"`
	Cover     string `json:"cover"`
	Rotate    int    `json:"rotate"`
	Likes     int    `json:"likes"`
	Dislikes  int    `json:"dislikes"`
	Listeners int    `json:"listeners"`
}

// StatisticCategory represents a statistics category.
type StatisticCategory struct {
	Key         string        `json:"key"`
	Description string        `json:"description,omitempty"`
	Icon        string        `json:"icon"`
	Tracks      []*TrackStats `json:"tracks"`
}

// StatisticsResponse represents the HTTP response for statistics.
type StatisticsResponse struct {
	Statistics []*StatisticCategory `json:"statistics"`
}
