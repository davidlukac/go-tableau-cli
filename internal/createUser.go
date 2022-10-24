package internal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/fastbill/go-httperrors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type CreateUserResponse struct {
	XMLName xml.Name            `xml:"tsResponse"`
	User    GetUserResponseUser `xml:"user"`
}

func (t Tableau) CreateUser(u User) (*User, error) {
	createUserURL := fmt.Sprintf("%s/sites/%s/users/", t.BaseURL, t.SiteID)

	log.Debugf("Searching for user on %s", createUserURL)

	requestBody := []byte(fmt.Sprintf(`
<tsRequest>
  <user name="%s" siteRole="%s" authSetting="%s" />
</tsRequest>
`, u.Username, u.Role, DefaultAuthSetting))

	req, err := http.NewRequest(http.MethodPost, createUserURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create create_user request: %w", err)
	}

	req.Header.Set("X-Tableau-Auth", t.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send create_user request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read 	response body: %w", err)
	}

	log.Debugf("response body: %s", string(body))
	log.Debugf("response code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		if resp.StatusCode == http.StatusUnauthorized {
			password := viper.GetString("tableau_password")
			return nil, fmt.Errorf("invalid credentials [%d] (%s:%s)", resp.StatusCode,
				viper.GetString("tableau_username"),
				fmt.Sprintf("%s*****%s", password[0:1], password[len(password)-1:]))
		} else if resp.StatusCode == http.StatusConflict {
			u.Exists = true
			return &u, fmt.Errorf("user %s already exists", u.Username)
		}
		return nil, httperrors.New(resp.StatusCode, fmt.Sprintf("failed to create user - server responded with status code: %d - %s", resp.StatusCode, string(body)))
	}

	var createUserResponse CreateUserResponse
	if err := xml.Unmarshal(body, &createUserResponse); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	log.Debugf("unmarshaled response: %v", createUserResponse.User)

	newUser := User{
		Username:    createUserResponse.User.Name,
		ID:          createUserResponse.User.ID,
		Role:        createUserResponse.User.SiteRole,
		AuthSetting: createUserResponse.User.AuthSetting,
		Exists:      true,
	}

	return &newUser, nil
}
