package user

type User struct {
	Id            string `json:"id"`
	Login         string `json:"login"`
	Name          string `json:"name"`
	Surname       string `json:"surname"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Pass          string `json:"pass"`
	Refresh_token string `json:"refresh_token"`
}
