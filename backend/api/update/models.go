package update

type result struct {
	Cacher string `json:"cacher"`
	Size   int    `json:"size"`
	Error  error  `json:"error"`
}
