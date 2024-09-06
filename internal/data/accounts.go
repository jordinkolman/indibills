package data

import (
	"database/sql"
	"errors"
)

type Account struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
	user_id int64
}

type AccountModel struct {
	DB *sql.DB
}

func (a AccountModel) Insert(account *Account) error {
	var type_id int64
	err := a.DB.QueryRow(`SELECT id FROM account_types WHERE type = $1`, account.Type).Scan(&type_id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errors.New("record not found")
		default:
			return err
		}
	}

	query := `
		INSERT INTO accounts (name, t., balance, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{account.Name, type_id, account.Balance, account.user_id}
	return a.DB.QueryRow(query, args...).Scan(&account.Id)

}

func (a AccountModel) GetAccountById(id int64) (*Account, error) {

	if id < 1 {
		return nil, errors.New("record not found")
	}

	query := `
	SELECT id, name, t.type, balance
	FROM accounts JOIN account_types t ON accounts.type_id = t.id
	WHERE id = $1`

	var account Account

	err := a.DB.QueryRow(query, id).Scan(
		&account.Id,
		&account.Name,
		&account.Type,
		&account.Balance,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}
	return &account, nil
}

func (a AccountModel) GetAccountsByUserId(id int64) ([]*Account, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	query := `
	SELECT id, name, t.type, balance
	FROM accounts JOIN account_types t ON account.type_id = t.id
	WHERE user_id = $1`

	rows, err := a.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	accounts := []*Account{}

	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.Id,
			&account.Name,
			&account.Type,
			&account.Balance,
		)
		if err != nil {
			return nil, err
		}

		account.user_id = id

		accounts = append(accounts, &account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}
