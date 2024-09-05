package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type User struct {
	Id        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UserResponse struct {
	User *User `json:"user"`
}

type UsersResponse struct {
	Users *[]User `json:"users"`
}

type UserListModel struct {
	Endpoint string
}

func (m *UserListModel) GetAll() (*[]User, error) {
	resp, err := http.Get(m.Endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var usersResp UsersResponse

	err = json.Unmarshal(data, &usersResp)
	if err != nil {
		return nil, err
	}
	return usersResp.Users, nil
}
