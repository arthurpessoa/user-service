package user

type User struct {
	Id    int
	Login string
	Name  string
}

type Repository interface {
	Find(id int) (*User)
	FindByLogin(login string) []*User
}

func NewUserRepository() Repository {
	return &userRepository{
		domains: make(map[int]*User),
	}
}

//Mocked implementation of the user repository
type userRepository struct {
	domains map[int]*User
}

func (*userRepository) Find(id int) (*User) {
	return nil
}

func (*userRepository) FindByLogin(login string) ([]*User) {
	return nil
}
