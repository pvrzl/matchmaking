package entity

type Status string

const (
	Up   Status = "UP"
	Down Status = "DOWN"
)

type HealthCheckApp struct {
	Name      string `json:"name"`
	BuildTag  string `json:"build_tag"`
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
}

type HealthcheckDatabase struct {
	Name   string `json:"name"`
	Status Status `json:"status"`
}

type HealthcheckComponents struct {
	Databases []HealthcheckDatabase `json:"databases"`
	Redis     Status                `json:"redis"`
}

type HealthcheckResponse struct {
	Status     Status                `json:"status"`
	App        HealthCheckApp        `json:"app"`
	Components HealthcheckComponents `json:"components"`
}
