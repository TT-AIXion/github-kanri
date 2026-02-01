package app

type syncResult struct {
	Repo   string `json:"repo"`
	Target string `json:"target"`
	Dest   string `json:"dest"`
}

type diffResult struct {
	Repo    string   `json:"repo"`
	Target  string   `json:"target"`
	Dest    string   `json:"dest"`
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Changed []string `json:"changed"`
}

type verifyResult struct {
	Repo   string `json:"repo"`
	Target string `json:"target"`
	Dest   string `json:"dest"`
	Match  bool   `json:"match"`
}
