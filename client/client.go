package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func MakeRequest(endpoint string, auth_token string, params url.Values) (*http.Response, error) {
	client := &http.Client{}
	req, _ := CreatePostRequest(endpoint, auth_token, params)

	resBod, _ := httputil.DumpRequest(req, true)
	fmt.Printf("Request Dump\n%s", string(resBod))

	resp, err := client.Do(req)

	return resp, err
}

func MakeGetRequest(endpoint string, auth_token string, params url.Values) (*http.Response, error) {
	client := &http.Client{}
	req, _ := CreateGetRequest(endpoint, auth_token, params)

	resBod, _ := httputil.DumpRequest(req, true)
	fmt.Printf("Request Dump\n%s", string(resBod))

	resp, err := client.Do(req)

	return resp, err
}

func Login(username, password string) (LoginResponse, error) {
	params := url.Values{}
	params.Add("username", username)
	params.Add("password", password)

	res, err := MakeRequest("/bq/login", STATIC_TOKEN, params)
	if res.StatusCode != 200 {
		return LoginResponse{}, errors.New("Login Unsuccessful")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Printf("JSON BODY:\n%s\n", body)
	PanicIfErr(err)

	var login LoginResponse
	err = json.Unmarshal(body, &login)

	PanicIfErr(err)

	return login, nil
}

// similar to a login, but does not authenticate
// simply retreives the snap response
func GetUpdates(username, auth_token string) (LoginResponse, error) {
	params := url.Values{}
	params.Add("username", username)

	res, err := MakeRequest("/bq/updates", auth_token, params)
	if res.StatusCode != 200 {
		return LoginResponse{}, errors.New("Update Failed")
	}

	PanicIfErr(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Printf("JSON BODY:\n%s\n", body)
	PanicIfErr(err)

	var login LoginResponse
	err = json.Unmarshal(body, &login)

	PanicIfErr(err)

	return login, nil
}

func Logout(username, auth_token string) bool {
	params := url.Values{}
	params.Add("username", username)
	params.Add("json", "{}")
	params.Add("events", "[]")

	res, err := MakeRequest("/ph/logout", auth_token, params)

	PanicIfErr(err)

	//logout returns no body, just looking for status 200
	if res.StatusCode == 200 {
		return true
	}

	return false
}

//returns the image/video bytes after being decrypted
func GetBlob(snap_id string, user User) []byte {
	params := url.Values{}
	params.Add("username", user.Username)
	params.Add("id", snap_id)

	res, err := MakeRequest("/ph/blob", user.AuthToken, params)

	PanicIfErr(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	//fmt.Printf("JSON BODY:\n%s", body)

	PanicIfErr(err)

	data := Decrypt(body)

	return data
}

func GetStories(username, auth_token string) StoryResponse {
	params := url.Values{}
	params.Add("username", username)

	res, err := MakeRequest("/bq/stories", auth_token, params)

	PanicIfErr(err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	PanicIfErr(err)

	var storyResponse StoryResponse
	err = json.Unmarshal(body, &storyResponse)

	PanicIfErr(err)

	return storyResponse
}

func GetStory(story_id, auth_token, media_key, media_iv string) []byte {
	params := url.Values{}
	params.Add("story_id", story_id)

	res, err := MakeGetRequest("/bq/story_blob", auth_token, params)
	PanicIfErr(err)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	PanicIfErr(err)

	fmt.Println("BOdy: ", body)

	data := DecryptStory(body, media_key, media_iv)
	return data
}
func PanicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
