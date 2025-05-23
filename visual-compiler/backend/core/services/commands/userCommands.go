package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"io"
)

func RegisterUser(username,email,password string) (string,error) {
	user_data := map[string]string{
		"username": username,
		"email": email,
		"password": password,
	}
	req,err := json.Marshal(user_data)
	if err != nil {
		fmt.Printf("error in core: could marshal")
		return "",fmt.Errorf("error in core: could marshal")
	}
	res, err := http.Post(
		"http://localhost:8080/api/register","application/json",
		bytes.NewBuffer(req),
	)
	if err != nil {
		return "",fmt.Errorf("error in core: could not send request to api:  %v",err)
	}

	defer res.Body.Close()

	if (res.StatusCode != http.StatusCreated) {
		bodyBytes, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("registration failed: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
        return "", fmt.Errorf("failed to read response body: %v", err)
    }
	var respMap map[string]string
    err = json.Unmarshal(bodyBytes, &respMap)
    if err != nil {
        return "", fmt.Errorf("failed to parse response JSON: %v", err)
    }

	if msg, ok := respMap["message"]; ok {
        return msg, nil
    }

    return "registration successful", nil
}