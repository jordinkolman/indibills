package data

import "database/sql"

type Models struct {
	Users UserModel
	Accounts AccountModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
	}
}
