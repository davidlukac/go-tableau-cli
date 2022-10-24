package internal

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

// DeleteUser - Remove user from site.
// https://help.tableau.com/current/api/rest_api/en-us/REST/rest_api_ref_users_and_groups.htm#remove_user_from_site
func (t Tableau) DeleteUser(userID, existingAssetsUserID string) (bool, error) {
	deleteUserURL := fmt.Sprintf("%s/sites/%s/users/%s", t.BaseURL, t.SiteID, userID)
	if existingAssetsUserID != "" {
		deleteUserURL = fmt.Sprintf("%s?mapAssetsTo=%s", deleteUserURL, existingAssetsUserID)
	}

	log.Debugf("Deleting user %s on URL %s", userID, deleteUserURL)

	req, err := http.NewRequest(http.MethodDelete, deleteUserURL, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create delete_user request: %w", err)
	}

	req.Header.Set("X-Tableau-Auth", t.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send delete_user request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debugf("response code: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusNoContent {
		if resp.StatusCode == http.StatusUnauthorized {
			password := viper.GetString("tableau_password")
			return false, fmt.Errorf("invalid credentials [%d] (%s:%s)", resp.StatusCode,
				viper.GetString("tableau_username"),
				fmt.Sprintf("%s*****%s", password[0:1], password[len(password)-1:]))
		}
		return false, fmt.Errorf("failed to delete user - server responded with status code: %d - %s", resp.StatusCode, string(body))
	}

	return true, nil
}
