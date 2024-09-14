package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
}

type UserModel struct {
	DB *sql.DB
}

func (u UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (email, password, firstName, lastName)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at`

	args := []interface{}{user.Email, user.Password, user.FirstName, user.LastName}
	return u.DB.QueryRow(query, args...).Scan(&user.Id, &user.CreatedAt)
}

func (u UserModel) Get(email string) (*User, error) {

	query := `
	SELECT id, email, password, firstName, lastName
	FROM users
	WHERE email = $1`

	var user User

	err := u.DB.QueryRow(query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}

	}
	fmt.Printf("Retrieved Successfully. \nId: %v\nEmail:%v\nFirstName:%v\n", user.Id, user.Email, user.FirstName)
	return &user, nil
}

func (u UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, firstName = $3, lastName = $4
		WHERE id = $6`
	args := []interface{}{user.Email, user.Password, user.FirstName, user.LastName}

	return u.DB.QueryRow(query, args...).Scan()
}

func (u UserModel) Delete(id int64) error {
	if id < 1 {
		return errors.New("record not found")
	}

	query := `
		DELETE FROM users
		WHERE id = $1`

	results, err := u.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (u UserModel) GetAll() ([]*User, error) {
	query := `
		SELECT *
		FROM users
		ORDER BY id`

	rows, err := u.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.CreatedAt,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
