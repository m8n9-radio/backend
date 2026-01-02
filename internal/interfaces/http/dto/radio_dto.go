package dto

// ListenerResponse represents listener count in HTTP response.
type ListenerResponse struct {
	Current int `json:"current"`
	Peak    int `json:"peak"`
}

// RadioResponse represents radio info in HTTP response.
type RadioResponse struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	StreamUrl   string           `json:"streamUrl"`
	Listener    ListenerResponse `json:"listener"`
}
