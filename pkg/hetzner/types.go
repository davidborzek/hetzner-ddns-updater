package hetzner

type (
	RecordType string

	Record struct {
		ID     string     `json:"id,omitempty"`
		ZoneID string     `json:"zone_id"`
		Type   RecordType `json:"type"`
		Name   string     `json:"name"`
		Value  string     `json:"value"`
		TTL    uint       `json:"ttl"`
	}
)

type (
	HetznerConsoleRRSet struct {
		Records []HetznerConsoleRecord `json:"records"`
	}

	HetznerConsoleRecord struct {
		Value string `json:"value"`
	}
)
