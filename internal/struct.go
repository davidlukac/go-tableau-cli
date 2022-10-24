package internal

import "encoding/xml"

type User struct {
	Username    string
	ID          string
	Role        string
	AuthSetting string
	Exists      bool
}

type Tableau struct {
	BaseURL string
	Token   string
	SiteID  string
}

type ErrorResponse struct {
	XLMName xml.Name           `xml:"tsResponse"`
	Error   ErrorResponseError `xml:"error"`
}

type ErrorResponseError struct {
	Code    string `xml:"code,attr"`
	Summary string `xml:"summary"`
	Detail  string `xml:"detail"`
}
