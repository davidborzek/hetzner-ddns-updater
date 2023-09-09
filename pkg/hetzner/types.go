package hetzner

type (
	Record struct {
		ID     string     `json:"id,omitempty"`
		ZoneID string     `json:"zone_id"`
		Type   RecordType `json:"type"`
		Name   string     `json:"name"`
		Value  string     `json:"value"`
		TTL    uint       `json:"ttl"`
	}

	RecordType string
)
