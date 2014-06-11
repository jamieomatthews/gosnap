package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jamieomatthews/gosnap/encryption"
)

const (
	//various constants
	URL string = "https://feelinsonice-hrd.appspot.com"

	SECRET              string = "iEk21fuwZApXlz93750dmW22pw389dPwOk"
	STATIC_TOKEN        string = "m198sOkJEn37DjqZ32lpRu76xmw288xSQ9"
	BLOB_ENCRYPTION_KEY string = "M02cnQ51Ji97vwT4"
	PATTERN             string = "0001110111101110001111010101111011010001001110011000110001000110"

	//media types
	IMAGE                        = 0
	VIDEO                        = 1
	VIDEO_NOAUDIO                = 2
	FRIEND_REQUEST               = 3
	FRIEND_REQUEST_IMAGE         = 4
	FRIEND_REQUEST_VIDEO         = 5
	FRIEND_REQUEST_VIDEO_NOAUDIO = 6

	//media states
	NONE       = -1
	SENT       = 0
	DELIVERED  = 1
	VIEWED     = 2
	SCREENSHOT = 3
)

func CreateRequestToken(token, timestamp string) string {
	hash := sha256.New()

	hash.Write([]byte(SECRET + token))
	md := hash.Sum(nil)
	first := hex.EncodeToString(md)
	hash = sha256.New()
	hash.Write([]byte(timestamp + SECRET))
	md = hash.Sum(nil)
	second := hex.EncodeToString(md)
	firstRune := []rune(first)
	secondRune := []rune(second)
	var buffer bytes.Buffer
	for i, ch := range PATTERN {
		if string(ch) == "0" {
			buffer.WriteString(string(firstRune[i]))
		} else {
			buffer.WriteString(string(secondRune[i]))
		}
	}
	return buffer.String()
}

func GetTimestamp() string {
	return strconv.Itoa(int(time.Now().Unix() * 1000))
}

func CreateGetRequest(endpoint, auth_token string, params url.Values) (*http.Request, error) {
	var s string = URL + endpoint
	now := GetTimestamp()
	params.Add("timestamp", now)
	params.Add("req_token", CreateRequestToken(auth_token, now))

	s = s + "?" + params.Encode()

	//create the request
	req, err := http.NewRequest("GET", s, nil)
	req.Header.Add("User-Agent", "Snapchat/6.1.2 (iPhone6,2; iOS 7.0.4; gzip)")

	return req, err
}

func CreatePostRequest(endpoint string, auth_token string, params url.Values) (*http.Request, error) {
	var s string = URL + endpoint
	now := GetTimestamp()
	params.Add("timestamp", now)
	params.Add("req_token", CreateRequestToken(auth_token, now))

	//create the request
	req, err := http.NewRequest("POST", s, bytes.NewBufferString(params.Encode()))
	req.Header.Add("User-Agent", "Snapchat/6.1.2 (iPhone6,2; iOS 7.0.4; gzip)")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	return req, err
}

func ConvertValuesToMap(params url.Values) map[string]string {
	vals := make(map[string]string)
	for key, val := range params {
		vals[key] = val[0]
	}
	return vals
}

func PKCS5Pad(data []byte) []byte {
	var blocksize uint8 = 16
	var padCount uint8 = 0
	padCount = blocksize - uint8(len(data)%int(blocksize))
	b := []byte{padCount}

	ba := bytes.Repeat(b, int(padCount))

	return append(data, ba...)
}

func Decrypt(ciphertext []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(BLOB_ENCRYPTION_KEY))
	mode := encryption.NewECBDecrypter(cipher)
	ciphertext = PKCS5Pad(ciphertext)
	mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}
func Encrypt(data []byte) []byte {
	cipher, _ := aes.NewCipher([]byte(BLOB_ENCRYPTION_KEY))
	cipher.Encrypt(data, PKCS5Pad(data))
	return data
}

func DecodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func EncodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func DecryptStory(data []byte, key, iv string) []byte {
	block, _ := aes.NewCipher([]byte(key))
	mode := cipher.NewCBCDecrypter(block, []byte(iv))
	ciphertext := PKCS5Pad(data)
	mode.CryptBlocks(ciphertext, ciphertext)
	return ciphertext
}

//leaving this out for now, as it depended on an external lib, and I'm not using it yet
func CreateMediaId(username string) string {
	//id, _ := uuid.NewV4()
	//return strings.ToUpper(username) + "~" + id.String()
	return ""
}

func IsVideo(data []byte) bool {
	//looking for byte values '\x00\x00'
	return len(data) > 1 && bytes.Equal(data[0:2], []byte{0, 0})
}

func IsImage(data []byte) bool {
	//looking for byte values '\xFF\xD8'
	//easiest way to do this is to just use the byte numbers
	return len(data) > 1 && bytes.Equal(data[0:2], []byte{255, 216})
}

//not yet implemented
func IsZip(data []byte) bool {
	return false
}
