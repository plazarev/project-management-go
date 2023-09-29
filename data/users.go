package data

type UsersStore struct{}

func (s *UsersStore) GetAll(ctx *DBContext) ([]User, error) {
	users := make([]User, 0)
	err := ctx.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, err
}
