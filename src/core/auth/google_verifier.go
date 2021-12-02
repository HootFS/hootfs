package google_verifer

import (
	"context"
	"encoding/json"
	"fmt"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
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

	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err == nil {
		log.Println("token: ", token)
	}
	username := g.get_username(ctx)
	isValid, err := g.VerifyAccessToken(token, username)
	if !isValid {
		return nil, fmt.Errorf("Error while logging in. \n %s", err)
	}
	return ctx, nil
}

func (g Google_verifier) get_username(ctx context.Context) string {
	val := metautils.ExtractIncoming(ctx).Get("username")
	return val
}

func (g Google_verifier) VerifyAccessToken(token string, username string) (bool, error) {
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
		return false, fmt.Errorf(string(body))
	}

	//for key, value := range result {
	//	// Each value is an interface{} type, that is type asserted as a string
	//	fmt.Println(key, value)
	//}

	if username == result["email"].(string) {
		log.Println("User successfully authenticated: ", username)
		return true, nil
	} else {
		print("Could not verify user.")
		return false, nil
	}
}
