package data

import "database/sql"

type Models struct {
	Users UserModel
	Accounts AccountStore
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{DB: db},
		Accounts: AccountStore{DB: db},
	}
}
