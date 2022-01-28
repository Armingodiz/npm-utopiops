package models

type AppDetail struct {
	State            State  `json:"state"`
	Name             string `json:"name"`
	Kind             string `json:"kind"`
	EnviironmentName string `json:"environmentName"`
	Status           string `json:"status"`
}

type State struct {
	Code string `json:"code"`
	Job  string `json:"job"`
}

type EnvrionmentDetail struct {
	State        State  `json:"state"`
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	ProviderName string `json:"providerName"`
	Status       string `json:"status"`
}
