package clouds

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"satellity/internal/configs"
)

var httpClient *http.Client

type recaptchaResp struct {
	Success bool    `json:"success"`
	Score   float64 `json:"score"`
}

// VerifyRecaptcha verify google recaptcha v3
func VerifyRecaptcha(ctx context.Context, recaptcha string) (bool, error) {
	config := configs.AppConfig
	if config.Environment == "test" {
		return true, nil
	}
	if len(config.Recaptcha.Secret) < 1 || len(config.Recaptcha.URL) < 1 {
		return true, nil
	}
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s?secret=%s&response=%s", config.Recaptcha.URL, config.Recaptcha.Secret, recaptcha), nil)
	if err != nil {
		return false, err
	}
	req.Close = true
	resp, err := httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	var captcha recaptchaResp
	err = json.Unmarshal(bytes, &captcha)
	if err != nil {
		return false, err
	}
	return captcha.Score > 0.6 && captcha.Success, nil
}
