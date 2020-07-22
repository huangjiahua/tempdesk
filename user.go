package tempdesk

const (
	NameAlreadyExists string = "name already exists"
)

type User struct {
	ID   int
	Name string
	Key  string
	Meta map[string]string
}

type UserService interface {
	User(name string) (user User, ok bool)
	CreateUser(user User) (err error)
	UpdateUser(user User) (err error)
}

type UserServiceError struct {
	Kind string
	Err  error
}

func (u *UserServiceError) Error() string {
	return u.Kind
}
