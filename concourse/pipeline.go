package concourse

type Pipeline struct {
	ID       int
	Name     string
	TeamName string `json:"team_name"`
}
