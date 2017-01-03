package util

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"gopkg.in/square/go-jose.v1/json"
)

const captchaURL = "https://www.google.com/recaptcha/api/siteverify"

type captchaResponse struct {
	Success bool
}

// CaptchaConfig struct used for the TOML configuration file
type CaptchaConfig struct {
	Enabled bool
	Public string
	Secret string
}

// VerifyCaptcha checks if the given captcha answer is valid
func VerifyCaptcha(answer string) (bool, error) {
	// Post form to google service
	resp, err := http.PostForm(captchaURL, url.Values{
		"secret": {
			Config.Captcha.Secret,
		},
		"response": {
			answer,
		},
	})

	// Check for errors
	if err != nil {
		return false, err
	}

	// Close body
	defer resp.Body.Close()

	// Read all content from response body
	body, err := ioutil.ReadAll(resp.Body)

	// Check for errors
	if err != nil {
		return false, err
	}

	captchaResp := captchaResponse{}

	// Unmarshal body to json struct
	if err := json.Unmarshal(body, &captchaResp); err != nil {

		return false, err
	}

	// Return success status
	return captchaResp.Success
}