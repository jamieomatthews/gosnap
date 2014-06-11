package client

// func TestRequest(t *testing.T) {
// 	fmt.Println("Testing Reqest\n\n")
// 	resp, err := MakeRequest("test", STATIC_TOKEN, url.Values{}, "GET")

// 	perror(err)

// 	body, _ := ioutil.ReadAll(resp.Body)
// 	fmt.Printf("%d : %s", resp.StatusCode, body)
// }
func perror(err error) {
	if err != nil {
		panic(err)
	}
}
