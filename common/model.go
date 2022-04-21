package common

// Git metadata info on local machine
type Git struct {
	// Url of git repo
	Url string `yaml:"url" json:"url"`
	// Branch of git repo
	Branch string `yaml:"branch" json:"branch"`
	// Tag of git repo
	Tag string `yaml:"tag" json:"tag"`
	// Commit info of git repo
	Commit *Commit `yaml:"commit" json:"commit"`
}

// Commit of git from local machine
type Commit struct {
	// Id of current commit
	Id string `yaml:"id" json:"id"`
	// Date of current commit
	Date string `yaml:"date" json:"date"`
	// IdAbbr is abbreviation of id of current commit
	IdAbbr string `yaml:"idAbbr" json:"idAbbr"`
	// Sub is subject of current commit
	Sub string `yaml:"sub" json:"sub"`
	// Committer of current commit
	Committer *Committer `yaml:"committer" json:"committer"`
}

// Committer info of current commit
type Committer struct {
	// Name of committer
	Name string `yaml:"name" json:"name"`
	// Email of committer
	Email string `yaml:"email" json:"email"`
}

type PkgInfo struct {
	Name    string
	Version string
}
