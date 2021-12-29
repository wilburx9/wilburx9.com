package update

type result struct {
	Cacher string  `json:"cacher,omitempty"`
	Size   int     `json:"size,omitempty"`
	Error  *errorV `json:"error,omitempty"`
}

type errorV struct {
	Message string `json:"message,omitempty"`
	Details error  `json:"details,omitempty"`
}
