package cmd

type Config struct {
	CoreUrl             string `json:"CORE_URL" yaml:"CORE_URL"`
	DMUrl               string `json:"DM_URL" yaml:"DM_URL"`
	LSMUrl              string `json:"LSM_URL" yaml:"LSM_URL"`
	UtopiopsOuthToken   string `json:"UTOPIOPS_OUTH_TOKEN" yaml:"UTOPIOPS_OUTH_TOKEN"`
	UtopiopsOuthIdToken string `json:"UTOPIOPS_OUTH_ID_TOKEN" yaml:"UTOPIOPS_OUTH_ID_TOKEN"`
	UtopiopsUserName    string `json:"UTOPIOPS_USERNAME" yaml:"UTOPIOPS_USERNAME"`
	UtopiopsPassword    string `json:"UTOPIOPS_PASSWORD" yaml:"UTOPIOPS_PASSWORD"`
	IdsUrl              string `json:"IDS_URL" yaml:"IDS_URL"`
	IdmUrl              string `json:"IDM_URL" yaml:"IDM_URL"`
}
