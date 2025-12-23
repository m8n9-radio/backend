package entity

type ListenerEntity struct {
	Current int `json:"current"`
	Peak    int `json:"peak"`
}

type RadioEntity struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	StreamUrl   string         `json:"streamUrl"`
	Listener    ListenerEntity `json:"listener"`
}
