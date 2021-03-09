package model

type Task struct {
	ID         string                 `json:"id,omitempty"`
	Name       string                 `json:"name,omitempty"`
	RoutingKey string                 `json:"routing_key,omitempty"` // NOTE if routing_key is empty, queue to main exchange
	Priority   int                    `json:"priority,omitempty"`
	Payload    map[string]interface{} `json:"payload"`
	Status     string                 `json:"status,omitempty"`
	Error      string                 `json:"error,omitempty"`
}
