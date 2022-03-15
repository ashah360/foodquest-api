package mock

import (
	"fmt"
	"github.com/ashah360/foodquest-api/internal/api/model"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type Mocker interface {
	InsertUsers(count int) error
}

type mocker struct {
	db *sqlx.DB
}

func (m *mocker) InsertUsers(count int) error {
	var users []*model.User

	for i := 0; i < count; i++ {
		users = append(users, NewUser())
	}

	var values []interface{}

	q := `INSERT INTO users (email, password, first_name, last_name, phone_number) values `

	nFields := 5

	for i, u := range users {
		values = append(values, u.Email, u.Password, u.FirstName, u.LastName, u.PhoneNumber)

		q += `(`
		for j := 0; j < nFields; j++ {
			q += `$` + strconv.Itoa(i*nFields+j+1) + `,`
		}
		q = q[:len(q)-1] + `),`
	}
	q = q[:len(q)-1]

	_, err := m.db.Exec(q, values...)

	fmt.Println(q)

	return err
}

func NewMocker(db *sqlx.DB) Mocker {
	return &mocker{
		db,
	}
}
