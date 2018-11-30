package nestedset

type (
	// Query ::
	Query struct {
		Get      string `json:"get"`       // all|child|parent
		NodeType int    `json:"node_type"` //
		NodeName string `json:"node_name"` //
	}
)
