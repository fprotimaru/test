package entity

type DepthType string

const (
	RateFactor DepthType = "factor"
	RateLimit  DepthType = "limit"
)

type Rate struct {
	Price  string    `json:"price"`
	Volume string    `json:"volume"`
	Amount string    `json:"amount"`
	Factor string    `json:"factor"`
	Type   DepthType `json:"type"`
}

type DepthAPIResponse struct {
	Timestamp int64  `json:"timestamp"`
	Asks      []Rate `json:"asks"`
	Bids      []Rate `json:"bids"`
}
