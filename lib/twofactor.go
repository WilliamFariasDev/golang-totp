package lib

//TwoFactor is for
type TwoFactor struct {
	KeySecret string `json:"key_secret"`
	User      User   `json:"user"`
}
