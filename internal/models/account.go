package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Account struct {
	Id      int64   `json:"id"`
	Name    string  `json:"name"`
	Type    string  `json:"type"`
	Balance float64 `json:"balance"`
	user_id int64 `json:"user_id"`
}

type AccountResponse struct {
	Account *Account `json:"account"`
}

type AccountsResponse struct {
	Accounts *[]Account `json:"accounts"`
}

type AccountModel struct {
	Endpoint string
}

type AccountListModel struct {
	Endpoint string
}

func (m *AccountListModel) GetAll(user_id int64) (*[]Account, error) {
	endpoint := fmt.Sprintf("%s/%v/accounts", m.Endpoint, user_id)
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.Body)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var accountsResp AccountsResponse


	err = json.Unmarshal(data, &accountsResp)
	if err != nil {
		return nil, err
	}

	return accountsResp.Accounts, nil
}
