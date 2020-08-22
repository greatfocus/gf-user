package services

// EmailRequest struct
type UserService struct {
	Host       string
	Port       string
	User       string
	Password   string
	From       string
	Subjects   []string
	Messages   []string
	Recipients []string
	Status     []bool
}

// entry function to create user
func (u *UserService) Create() {
}
