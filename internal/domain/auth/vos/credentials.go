package vos

type Credentials struct {
	ClientID     string
	ClientSecret string
	Scope        string
	Extra        map[string]string
}