package nestedset

type (
	// Query ::
	Query struct {
		Get      string `json:"get"`       // all|child|parent
		NodeType int32  `json:"node_type"` //
		NodeName string `json:"node_name"` //
	}
)
