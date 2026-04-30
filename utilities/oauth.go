package utilities

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"forum/models"
)

// func generateState() string {
// 	b := make([]byte, 16)
// 	rand.Read(b)
// 	return base64.URLEncoding.EncodeToString(b)
// }

// In production, store states in Redis/DB with expiry.
// This in-memory map is fine for a single instance.
// var stateStore = map[string]bool{}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	IDToken     string `json:"id_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`

	// Scope       string `json:"scope"`
	Error string `json:"error"`
}

// type UserInfo struct {
// 	ID      int    `json:"id"` // check for google
// 	Name    string `json:"name"`
// 	Email   string `json:"email"`
// 	Picture string `json:"picture"`    // gmail picture: sometimes cannot be loaded!
// 	Avatar  string `json:"avatar_url"` // github avatar
// }

var redirectURI = "http://localhost:8080/auth/google/callback"

func ExchangeCode(provider, tokenURL, client_id, client_secret, code string) (*TokenResponse, error) {
	var resp *http.Response
	var err error

	if provider == "google" {
		data := url.Values{}
		data.Set("code", code)
		data.Set("client_id", client_id)
		data.Set("client_secret", client_secret)

		// only for google ?!
		data.Set("redirect_uri", redirectURI)

		data.Set("grant_type", "authorization_code")

		resp, err = http.PostForm(tokenURL, data)
		if err != nil {
			return nil, err
		}
		// defer resp.Body.Close()
	} else {
		// 2. for github !
		// resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded",
		// 	strings.NewReader(body.Encode()))
		// if err != nil {
		// 	return nil, err
		// }
		// defer resp.Body.Close()
		payload := map[string]string{
			"client_id":     client_id,
			"client_secret": client_secret,
			"code":          code,
		}

		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return nil, err
		}

		// IMPORTANT: Ask for JSON response
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
	}

	// ===================================

	// if resp.StatusCode != http.StatusOK {
	// 	raw, _ := io.ReadAll(resp.Body)
	// 	return nil, fmt.Errorf("status %d: %s", resp.StatusCode, raw)
	// }

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	tokenData := TokenResponse{}
	json.Unmarshal(body, &tokenData)
	return &tokenData, nil
	// var t TokenResponse
	// if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
	// 	return nil, err
	// }
	// return &t, nil
}

func FetchUserInfo(userInfoURL, accessToken string) (*models.User, error) {
	req, _ := http.NewRequest("GET", userInfoURL, nil)
	// for github: The /user endpoint returns email ONLY IF it is public
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// IMPORTANT: Ask for JSON response
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{} // or: http.DefaultClient
	userResp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer userResp.Body.Close()

	userBody, _ := io.ReadAll(userResp.Body)

	var user models.User
	json.Unmarshal(userBody, &user)
	return &user, nil
}

func FetchGithubUserEmail(accessToken string) (string, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json") // use it everywhere, why ?!

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	var primary string
	for _, e := range emails {
		if e["primary"].(bool) {
			primary = e["email"].(string)
			break
		}
	}

	return primary, nil
}
