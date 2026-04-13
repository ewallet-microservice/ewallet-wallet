package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mhasnanr/ewallet-wallet/bootstrap"
	"github.com/mhasnanr/ewallet-wallet/constants"
)

type ExternalUserAPI struct{}

type ValidateUserResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	Email    string `json:"email"`
}

type UserValidationResponse struct {
	Message string               `json:"message"`
	Data    ValidateUserResponse `json:"data"`
}

func (e *ExternalUserAPI) ValidateToken(accessToken string) (ValidateUserResponse, error) {
	var (
		err    error
		user   = ValidateUserResponse{}
		client = &http.Client{}
	)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/v1/token/validate", bootstrap.GetEnv("USER_API_BASE_URL", "")), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)

	if err != nil {
		return user, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return user, err
	}

	if resp.StatusCode != http.StatusOK {
		errorStr := fmt.Errorf("user service error: status %d: %s", resp.StatusCode, string(body))
		fmt.Println(errorStr)
		return user, constants.ErrorUnauthorized
	}

	var wrapper UserValidationResponse
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return user, constants.ErrorFailedToParseUser
	}

	return wrapper.Data, nil
}
