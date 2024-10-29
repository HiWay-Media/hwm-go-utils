package nomad

type ScaleRequest struct {
	Count  int               `json:"Count,omitempty"`
	Target ScaleGroupRequest `json:"Target"`
}

type ScaleGroupRequest struct {
	Group string `json:"Group,omitempty"`
}

type RunJobRequest struct {
	Job    JobDefinition `json:"Job"`
	Format string        `json:"Format,omitempty"`
}
