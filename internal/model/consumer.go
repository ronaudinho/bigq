package model

type Consumer struct {
	Tag        string `json:"tag"`
	RoutingKey string `json:"routing_key"`
}
