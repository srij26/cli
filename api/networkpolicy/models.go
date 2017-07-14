package networkpolicy

type Policy struct {
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
}

type Source struct {
	ID  string `json:"id"`
	Tag string `json:"tag,omitempty"`
}

type Destination struct {
	ID       string `json:"id"`
	Tag      string `json:"tag,omitempty"`
	Protocol string `json:"protocol"`
	Port     int    `json:"port,omitempty"`
	Ports    Ports  `json:"ports"`
}

type Ports struct {
	Start int `json:"start"`
	End   int `json:"end"`
}
