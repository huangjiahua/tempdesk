package tempdesk

type User struct {
	ID   int
	Name string
	Key  string
}

type UserService interface {
	User(name string) (user User, ok bool)
	CreateUser(user User) (err error)
}
