package api

import (
	"encoding/json"
	"errors"

	"github.com/Adeithe/go-twitch/api/request"
)

type loginData struct {
	ClientID string  `json:"client_id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Undelete bool    `json:"undelete_user"`
	Captcha  captcha `json:"captcha,omitempty"`
}

type captcha struct {
	Proof string `json:"proof,omitempty"`
}

// TwitchLogin contains data for a Twitch login request
type TwitchLogin struct {
	Username      string
	password      string
	CaptchaProof  string `json:"captcha_proof,omitempty"`
	ObscuredEmail string `json:"obscured_email,omitempty"`
	Error         string `json:"error,omitempty"`
	ErrorCode     int    `json:"error_code"`
	ErrorShort    string `json:"error_description,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	RedirectPath  string `json:"redirect_path,omitempty"`
}

// ITwitchLogin interface containing methods for the Twitch Login
type ITwitchLogin interface {
	Verify(string) error
	ToBearer() (*Client, error)
	GetError() string
	GetErrorCode() int
	GetAccessToken() string
}

var _ ITwitchLogin = &TwitchLogin{}

// Verify may need to be called if 2FA is enabled
func (login *TwitchLogin) Verify(code string) error {
	var marshalFunc func(string, loginData) ([]byte, error)

	switch login.ErrorCode {
	// unknown
	case -1:
		marshalFunc = marshalTwitchguardData

	// authy
	case 3011:
		marshalFunc = marshalAuthyData

	// twitchguard
	case 3022:
		marshalFunc = marshalTwitchguardData

	default:
		return errors.New("verification code not required")
	}

	body, err := marshalFunc(code, loginData{
		ClientID: Official.ID,
		Username: login.Username,
		Password: login.password,
		Captcha:  captcha{login.CaptchaProof},
	})
	if err != nil {
		return err
	}

	req := request.New("POST", "https://passport.twitch.tv", "login")
	req.Headers["Content-Type"] = "application/json"
	req.Body = body

	res, err := req.Do()
	if err != nil {
		return err
	}
	if err := json.Unmarshal(res.Body, &login); err != nil {
		return err
	}
	if len(login.AccessToken) > 0 {
		login.ErrorCode = 0
		login.Error = ""
		login.ErrorShort = "login success"
		login.CaptchaProof = ""
	}

	return nil
}

// ToBearer creates a API client using the TwitchLogin.
func (login TwitchLogin) ToBearer() (*Client, error) {
	if len(login.AccessToken) <= 0 {
		return nil, errors.New(login.ErrorShort)
	}
	return Official.NewBearer(login.AccessToken), nil
}

// GetError returns the current login error in readable text.
// This is the same text that would appear on the Twitch website.
func (login TwitchLogin) GetError() string {
	return login.Error
}

// GetErrorCode returns the current login error code.
// Returns -1 if the login attempt failed completely and 0 if the login was a success
func (login TwitchLogin) GetErrorCode() int {
	return login.ErrorCode
}

// GetAccessToken returns the users Bearer token.
// Returns an empty string if the login process is not yet finished.
func (login TwitchLogin) GetAccessToken() string {
	return login.AccessToken
}

func marshalTwitchguardData(code string, ld loginData) ([]byte, error) {
	return json.Marshal(struct {
		loginData
		Verification string `json:"twitchguard_code"`
	}{
		loginData:    ld,
		Verification: code,
	})
}

func marshalAuthyData(code string, ld loginData) ([]byte, error) {
	return json.Marshal(struct {
		loginData
		Verification string `json:"authy_token"`
	}{
		loginData:    ld,
		Verification: code,
	})
}
