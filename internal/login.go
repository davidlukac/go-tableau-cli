package internal

import (
	"bytes"
	"encoding/xml"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type LoginResponse struct {
	XMLName     xml.Name                 `xml:"tsResponse"`
	Credentials LoginResponseCredentials `xml:"credentials"`
}

type LoginResponseCredentials struct {
	Token string            `xml:"token,attr"`
	Site  LoginResponseSite `xml:"site"`
	User  LoginResponseUser `xml:"user"`
}

type LoginResponseSite struct {
	ID string `xml:"id,attr"`
}

type LoginResponseUser struct {
	ID string `xml:"id,attr"`
}

func Login(baseURL, username, password string) (token, siteId string, err error) {
	loginURL := fmt.Sprintf("%s/auth/signin", baseURL)

	var payload = []byte(fmt.Sprintf(`
<tsRequest>
  <credentials name="%s" password="%s" >
    <site contentUrl="" />
  </credentials>
</tsRequest>
`, username, password))

	req, err := http.NewRequest(http.MethodPost, loginURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", "", fmt.Errorf("failed to create login request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to send login request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debugln(string(body))

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to log in - server responded with status code: %d - %s", resp.StatusCode, string(body))
	}

	var loginResponse LoginResponse
	if err := xml.Unmarshal(body, &loginResponse); err != nil {
		return "", "", fmt.Errorf("unable to unmarshal response body: %w", err)
	}

	return loginResponse.Credentials.Token, loginResponse.Credentials.Site.ID, nil
}
