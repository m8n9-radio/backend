package radio

// RadioInfo represents radio stream information.
type RadioInfo struct {
	name         string
	description  string
	streamURL    string
	listeners    int
	listenerPeak int
}

// NewRadioInfo creates a new RadioInfo value object.
func NewRadioInfo(name, description, streamURL string, listeners, listenerPeak int) RadioInfo {
	return RadioInfo{
		name:         name,
		description:  description,
		streamURL:    streamURL,
		listeners:    listeners,
		listenerPeak: listenerPeak,
	}
}

// Name returns the radio name.
func (r RadioInfo) Name() string {
	return r.name
}

// Description returns the radio description.
func (r RadioInfo) Description() string {
	return r.description
}

// StreamURL returns the stream URL.
func (r RadioInfo) StreamURL() string {
	return r.streamURL
}

// Listeners returns the current listener count.
func (r RadioInfo) Listeners() int {
	return r.listeners
}

// ListenerPeak returns the peak listener count.
func (r RadioInfo) ListenerPeak() int {
	return r.listenerPeak
}

// IsEmpty returns true if the RadioInfo is empty.
func (r RadioInfo) IsEmpty() bool {
	return r.name == "" && r.streamURL == ""
}
