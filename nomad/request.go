package nomad

type ScaleRequest struct {
	Count  int               `json:"Count"`
	Target ScaleGroupRequest `json:"Target"`
}

type ScaleGroupRequest struct {
	Group string `json:"Group"`
}

type RunJobRequest struct {
	Job    JobDefinition `json:"Job"`
	Format string        `json:"Format"`
}
