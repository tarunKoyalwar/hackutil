package burpconfig

type BurpConfig struct {
	Target Target `json:"target"`
}

type Target struct {
	Scope Scope `json:"scope"`
}

type Scope struct {
	AdvancedMode bool    `json:"advanced_mode"`
	Exclude      []Rule  `json:"exclude"`
	Include      []Rule  `json:"include"`
}

type Rule struct {
	Enabled  bool   `json:"enabled"`
	File     string `json:"file"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
} 