package data

import (
	"database/sql"
	"testing"
	"time"
)

type PostgressRepositoryTest struct {
	Conn *sql.DB
}

func NewPostgressRepositoryTest(dbPool *sql.DB) *PostgressRepositoryTest {
	return &PostgressRepositoryTest{
		Conn: dbPool,
	}
}

func (r *PostgressRepositoryTest) GetAll() ([]*User, error) {
	users := []*User{}

	return users, nil
}

func (r *PostgressRepositoryTest) GetByEmail(email string) (*User, error) {
	user := User{
		ID:        1,
		Email:     "bla@bla.com",
		FirstName: "Name",
		LastName:  "Last",
		Password:  "bla",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

// GetOne returns one user by id
func (r *PostgressRepositoryTest) GetOne(id int) (*User, error) {
	user := User{
		ID:        id,
		Email:     "bla@bla.com",
		FirstName: "Name",
		LastName:  "Last",
		Password:  "bla",
		Active:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &user, nil
}

func (r *PostgressRepositoryTest) Update(u User) error {
	return nil
}

func (r *PostgressRepositoryTest) DeleteByID(id int) error {
	return nil
}

func (r *PostgressRepositoryTest) Insert(user User) (int, error) {
	return 2, nil
}

func (r *PostgressRepositoryTest) ResetPassword(password string, u User) error {
	return nil
}

func (r *PostgressRepositoryTest) PasswordMatches(plainText string, u User) (bool, error) {
	return true, nil
}

func Test_models(t *testing.T) {

}
