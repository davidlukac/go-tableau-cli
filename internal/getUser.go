package internal

import (
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type GetUserResponse struct {
	XLMName    xml.Name              `xml:"tsResponse"`
	Pagination Pagination            `xml:"pagination"`
	Users      []GetUserResponseUser `xml:"users>user"`
}

type GetUserResponseUser struct {
	ID          string `xml:"id,attr"`
	Name        string `xml:"name,attr"`
	SiteRole    string `xml:"siteRole,attr"`
	AuthSetting string `xml:"authSetting,attr"`
}

// Pagination : <pagination pageNumber="1" pageSize="100" totalAvailable="341"/>
type Pagination struct {
	PageNumber     int `xml:"pageNumber,attr"`
	PageSize       int `xml:"pageSize,attr"`
	TotalAvailable int `xml:"totalAvailable,attr"`
}

// GetUser - Get user object from the server.
// Returns non-nil err object if
// - fails to construct request or client
// - fails to parse response
// - fails to log in or other non-success code is returned from the server
// - more than one user is returned in the result (match is to exact username)
// Returns empty User struct with Exists set to false if no matches were found.
// Returns User object on success.
// API Doc: https://help.tableau.com/current/api/rest_api/en-us/REST/rest_api_ref_users_and_groups.htm#get_users_on_site
// API Endpoint: GET /api/api-version/sites/site-id/users?filter=filter-expression
func (t Tableau) GetUser(username string) (*User, error) {
	searchUserURL := fmt.Sprintf("%s/sites/%s/users/?filter=name:eq:%s", t.BaseURL, t.SiteID, username)

	log.Debugf("Searching for user on %s", searchUserURL)

	req, err := http.NewRequest(http.MethodGet, searchUserURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get_user request: %w", err)
	}

	req.Header.Set("X-Tableau-Auth", t.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send get_user request: %w", err)
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
		return nil, fmt.Errorf("failed to get user - server responded with status code: %d - %s", resp.StatusCode, string(body))
	}

	var getUserResponse GetUserResponse
	if err := xml.Unmarshal(body, &getUserResponse); err != nil {
		return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	if len(getUserResponse.Users) == 0 {
		return &User{Exists: false}, nil
	}
	if len(getUserResponse.Users) > 1 {
		return &User{Exists: false}, fmt.Errorf("ambiguous result - more than one user returned")
	}

	user := getUserResponse.Users[0]
	return &User{
		Exists:   true,
		Username: user.Name,
		ID:       user.ID,
		Role:     user.SiteRole,
	}, nil
}

// GetUsers returns list of all users in given site.
// API Doc: https://help.tableau.com/current/api/rest_api/en-us/REST/rest_api_ref_users_and_groups.htm#get_users_on_site
// API Endpoint: GET /api/api-version/sites/site-id/users
func (t Tableau) GetUsers() ([]*User, error) {
	done := false
	pageNumber := 1
	pageSize := 100
	res := make([]*User, 0)

	for done == false {
		url := fmt.Sprintf("%s/sites/%s/users?pageSize=%d&pageNumber=%d", t.BaseURL, t.SiteID, pageSize, pageNumber)

		log.Debugf("Fetching %d users/page %d from %s", pageSize, pageNumber, url)

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create get_users request: %w", err)
		}

		req.Header.Set("X-Tableau-Auth", t.Token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send get_users request: %w", err)
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
			return nil, fmt.Errorf("failed to get user - server responded with status code: %d - %s", resp.StatusCode, string(body))
		}

		var getUserResponse GetUserResponse
		if err := xml.Unmarshal(body, &getUserResponse); err != nil {
			return nil, fmt.Errorf("unable to unmarshal response body: %w", err)
		}

		users := getUserResponse.Users

		if len(users) == 0 {
			log.Info("No users were found")
			return []*User{}, nil
		}

		log.Debugf("Server returned %d users.", len(users))

		for _, user := range users {
			res = append(res, &User{
				Exists:   true,
				Username: user.Name,
				ID:       user.ID,
				Role:     user.SiteRole,
			})
		}

		done = len(res) >= getUserResponse.Pagination.TotalAvailable
		pageNumber++
	}

	return res, nil
}
