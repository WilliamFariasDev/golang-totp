package lib

//User is for
type User struct {
	Email     string `json:"email"`
	UUID      string `json:"uuid"`
	TwoFactor int    `json:"two_factor"`
	Key       string `json:"key"`
}
