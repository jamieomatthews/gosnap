package client

type User struct {
	Username      string `form:"username" binding:"required"`
	Password      string `form:"password" binding:"required"`
	AuthToken     string
	FriendStories StoryResponse //store the last gotten story response
}

//////////////////////////////////
//Login Structs
//////////////////////////////////
type LoginResponse struct {
	Bests     []string `json:"bests"`
	Score     int      `json:"score"`
	Snaps     []Snap   `json:"snaps"`
	Friends   []Friend `json:"friends"`
	AuthToken string   `json:"auth_token"`
	Username  string   `json:"username"`
}

type Snap struct {
	SnapId        string `json:"id"`
	ScreenName    string `json:"sn"`
	RecipientName string `json:"rp"`
	MediaType     int    `json:"m"`
	MediaState    int    `json:"st"`
	Unopened      uint64 `json:"t,omitempty"`
}

type Friend struct {
	Name        string `json:"name"`
	DisplayName string `json:"display"`
}

func (s Snap) IsIncoming() bool {
	return s.RecipientName == ""
}

func (s Snap) IsUnopened() bool {
	return s.Unopened != 0
}

func (s Snap) IsImage() bool {
	return s.MediaType == IMAGE
}

func (s Snap) IsVideo() bool {
	return s.MediaType == VIDEO
}

//////////////////////////////////
//Story Structs
//////////////////////////////////

type StoryResponse struct {
	Friends []FriendStoryDict `json:"friend_stories"`
}

//User mapped to a list of that users available stories
type FriendStoryDict struct {
	Username      string        `json:"username"`
	FriendStories []FriendStory `json:"stories"`
}

//individual story object
type FriendStory struct {
	Viewed  bool         `json:"viewed"`
	Stories StoryContent `json:"story"`
}

//inner content for a story
type StoryContent struct {
	Id           string `json:"id"`
	Username     string `json:"username"`
	MediaId      string `json:"media_id"`
	MediaKey     string `json:"media_key"`
	MediaIv      string `json:"media_iv"`
	ThumbnailIv  string `json:"thumbnail_iv"`
	MediaType    int    `json:"media_type"`
	MediaUrl     string `json:"media_url"`
	ThumbnailUrl string `json:"thumbnail_url"`
}
