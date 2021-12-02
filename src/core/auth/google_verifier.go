package google_verifer

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Google_verifier struct {
	Url string
}

func New(url string) Google_verifier {
	googleV := Google_verifier{url}
	return googleV
}

func (g Google_verifier) Authenticate(ctx context.Context) (context.Context, error) {
	print("Triggered Auth")
	log.Println("TRIGGERED")
	return ctx, nil
}

func (g Google_verifier) VerifyAccessToken(token string, username string) bool {
	//Reference: https://stackoverflow.com/questions/51452148/how-can-i-make-a-request-with-a-bearer-token-in-go
	var bearer = "Bearer " + token

	// Create a new request using http
	req, err := http.NewRequest("GET", g.Url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		print("Error: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	_, err = json.Marshal(&body)

	if err != nil {
		log.Println("Could not parse the json: ", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)

	if err != nil {
		log.Println("Could not unmarshal the json: ", err)
	}

	//for key, value := range result {
	//	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(key, value)
	//}

	if username == result["email"].(string) {
		log.Println("User successfully authenticated: ", username)
		return true
	} else {
		print("Could not verify user.")
		return false
	}
}
