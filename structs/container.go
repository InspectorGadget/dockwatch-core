package structs

type Container struct {
	Name   string          `json:"name"`
	State  string          `json:"state"`
	Status string          `json:"status"`
	Stat   ContainerMetric `json:"stats"`
}

type ContainerMetric struct {
	CPU      string `json:"cpu"`
	MemUsage string `json:"mem_usage"`
	MemLimit string `json:"mem_limit"`
	MemPerc  string `json:"mem_percent"`
}
