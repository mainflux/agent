package sdk

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mainflux/mainflux/pkg/errors"
)

// Token is used for authentication purposes.
// It contains AccessToken, RefreshToken and AccessExpiry.
type Token struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessType   string `json:"access_type,omitempty"`
}

func (sdk mfSDK) CreateToken(user User) (Token, errors.SDKError) {
	var treq = tokenReq{
		Identity: user.Credentials.Identity,
		Secret:   user.Credentials.Secret,
	}
	data, err := json.Marshal(treq)
	if err != nil {
		return Token{}, errors.NewSDKError(err)
	}

	url := fmt.Sprintf("%s/%s/%s", sdk.usersURL, usersEndpoint, issueTokenEndpoint)

	_, body, sdkerr := sdk.processRequest(http.MethodPost, url, "", string(CTJSON), data, http.StatusCreated)
	if sdkerr != nil {
		return Token{}, sdkerr
	}
	var token Token
	if err := json.Unmarshal(body, &token); err != nil {
		return Token{}, errors.NewSDKError(err)
	}

	return token, nil
}

func (sdk mfSDK) RefreshToken(token string) (Token, errors.SDKError) {
	url := fmt.Sprintf("%s/%s/%s", sdk.usersURL, usersEndpoint, refreshTokenEndpoint)

	_, body, sdkerr := sdk.processRequest(http.MethodPost, url, token, string(CTJSON), []byte{}, http.StatusCreated)
	if sdkerr != nil {
		return Token{}, sdkerr
	}

	var t = Token{}
	if err := json.Unmarshal(body, &t); err != nil {
		return Token{}, errors.NewSDKError(err)
	}

	return t, nil
}
