package supabase

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrUserExists    = errors.New("User already exists")
	ErrWrongPassword = errors.New("Wrong password")
)

type Client interface {
	SignupUser(email, password string) (CreateUserResponse, error)
	SigninUser(email, password string) (SignInUserResponse, error)
}

type client struct {
	url    string
	secret string
}

func NewClient(url, secret string) Client {
	return &client{
		url:    url,
		secret: secret,
	}
}

func (c *client) SigninUser(email, password string) (SignInUserResponse, error) {
	url := c.url + "/auth/v1/token?grant_type=password"

	reqPl := CreateUserRequest{
		Email:    email,
		Password: password,
	}

	pl, err := json.Marshal(reqPl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pl))
	if err != nil {
		return SignInUserResponse{}, err
	}

	req.Header.Add("apikey", c.secret)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SignInUserResponse{}, err
	}

	defer func() {
		res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return SignInUserResponse{}, err
	}

	if res.StatusCode != 200 {
		ret := LoginErrorResponse{}
		if err := json.Unmarshal(body, &ret); err != nil {
			return SignInUserResponse{}, err
		}

		if ret.Error == "invalid_grant" {
			return SignInUserResponse{}, ErrWrongPassword
		}

		return SignInUserResponse{}, errors.New(ret.ErrorDescription)
	}

	ret := SignInUserResponse{}
	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}

func (c *client) SignupUser(email, password string) (CreateUserResponse, error) {
	url := c.url + "/auth/v1/signup"

	reqPl := CreateUserRequest{
		Email:    email,
		Password: password,
	}

	pl, err := json.Marshal(reqPl)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pl))
	if err != nil {
		return CreateUserResponse{}, err
	}

	req.Header.Add("apikey", c.secret)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CreateUserResponse{}, err
	}

	defer func() {
		res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CreateUserResponse{}, err
	}

	if res.StatusCode != 200 {
		ret := ErrorResponse{}
		if err := json.Unmarshal(body, &ret); err != nil {
			return CreateUserResponse{}, err
		}

		if ret.Msg == "User already registered" {
			return CreateUserResponse{}, ErrUserExists
		}

		return CreateUserResponse{}, errors.New(ret.Msg)
	}

	ret := CreateUserResponse{}
	if err := json.Unmarshal(body, &ret); err != nil {
		return ret, err
	}

	return ret, nil
}
