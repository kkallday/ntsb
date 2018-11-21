package concourse

type Build struct {
	ID       int
	Name     string
	Status   string
	JobName  string `json:"job_name"`
	TeamName string `json:"team_name"`
}
