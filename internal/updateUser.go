package internal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

type UpdateUserResponse struct {
	XLMName xml.Name               `xml:"tsResponse"`
	User    UpdateUserResponseUser `xml:"user"`
}

type UpdateUserResponseUser struct {
	Name        string `xml:"name,attr"`
	FullName    string `xml:"fullName,attr"`
	Email       string `xml:"email,attr"`
	SiteRole    string `xml:"siteRole,attr"`
	AuthSetting string `xml:"authSetting,attr"`
}

// UpdateUserSiteRole updates existing user site role to provided role.
// Returns non-nil error object if there was an error getting, updating or verifying the update.
// Returns User object with exists set to false if user was not found.
// API Doc: https://help.tableau.com/current/api/rest_api/en-us/REST/rest_api_ref_users_and_groups.htm#update_user
// API Endpoint: PUT /api/api-version/sites/site-id/users/user-id
func (t Tableau) UpdateUserSiteRole(username, siteRole string) (*User, error) {
	user, err := t.GetUser(username)
	if err != nil {
		return nil, err
	}
	if false == user.Exists {
		return user, nil
	}
	if strings.ToLower(siteRole) == strings.ToLower(user.Role) {
		log.Infof("User %s already has role %s assigned.", user.Username, user.Role)
		return user, nil
	}

	updateUserURL := fmt.Sprintf("%s/sites/%s/users/%s", t.BaseURL, t.SiteID, user.ID)

	var payload = []byte(fmt.Sprintf(`
<tsRequest>
  <user siteRole="%s" />
</tsRequest>
`, siteRole))

	log.Debugf("Updating user on URL %s", updateUserURL)

	req, err := http.NewRequest(http.MethodPut, updateUserURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create update_user request: %w", err)
	}

	req.Header.Set("X-Tableau-Auth", t.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send update_user request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debugf("response body: %s", string(body))
	log.Debugf("response code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusUnauthorized {
			password := viper.GetString("tableau_password")
			return nil, fmt.Errorf("invalid credentials [%d] (%s:%s)", resp.StatusCode,
				viper.GetString("tableau_username"),
				fmt.Sprintf("%s*****%s", password[0:1], password[len(password)-1:]))
		}

		var updateUserErrorResponse ErrorResponse
		if err := xml.Unmarshal(body, &updateUserErrorResponse); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
		}

		return nil, fmt.Errorf("failed to update user - server responded with status code: %d, Code: %s, "+
			"Summary: %s, Detail: %s", resp.StatusCode, updateUserErrorResponse.Error.Code,
			updateUserErrorResponse.Error.Summary, updateUserErrorResponse.Error.Detail)
	}

	var updateUserResponse UpdateUserResponse
	if err := xml.Unmarshal(body, &updateUserResponse); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	user, err = t.GetUser(username)
	if err != nil {
		return nil, err
	}
	if strings.ToLower(user.Role) != strings.ToLower(updateUserResponse.User.SiteRole) ||
		strings.ToLower(user.Role) != strings.ToLower(siteRole) {
		return nil, fmt.Errorf("something with wrong - updated and requested roles don't match: %s - %s - %s", siteRole, updateUserResponse.User.SiteRole, user.Role)
	}

	return user, nil
}
