package main

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/jamieomatthews/gosnap/client"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

//for now, we will store users in memory
//this should probably be expanded to use a database
var users map[string]client.User = make(map[string]client.User)

func main() {
	m := martini.Classic()

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("gosnap", store))
	m.Use(render.Renderer())

	//render out the inital page
	m.Get("/", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	//////////////////////////
	//API
	//////////////////////////

	//attempt to reload your snaps by checking the session variable
	m.Get("/api/reload", authorize, func(currentUser client.User, r render.Render) {
		var loginRequest = client.Login(currentUser.Username, currentUser.Password)
		currentUser.AuthToken = loginRequest.AuthToken
		users[currentUser.Username] = currentUser
		r.HTML(200, "snaps", loginRequest)
	})

	//log the user in, and store them in memory
	m.Post("/api/login", binding.Bind(client.User{}), func(user client.User, r render.Render, session sessions.Session) {
		var loginRequest = client.Login(user.Username, user.Password)

		//store the user in memory
		session.Set("user_id", user.Username)
		user.AuthToken = loginRequest.AuthToken
		users[user.Username] = user

		r.HTML(200, "snaps", loginRequest)
	})

	m.Get("/api/logout", authorize, func(user client.User, r render.Render, session sessions.Session) {
		didLogout := client.Logout(user.Username, user.AuthToken)
		if didLogout {
			//Remove the user from session, and from memory
			session.Delete("user_id")
			delete(users, user.Username)

			r.JSON(200, map[string]string{"status": "ok"})
			return
		}
		r.JSON(400, map[string]string{"status": "error"})
	})

	//return the users stories
	m.Get("/api/stories", authorize, func(user client.User, r render.Render) {
		storyResponse := client.GetStories(user.Username, user.AuthToken)
		//store the story response for later retrieval
		user.FriendStories = storyResponse
		users[user.Username] = user
		r.HTML(200, "stories", storyResponse)
	})

	//given the user is logged in, return the snap by ID
	m.Post("/api/snap", authorize, func(res http.ResponseWriter, req *http.Request, r render.Render, session sessions.Session, currentUser client.User) {
		id := req.FormValue("snap_id")
		data := client.GetBlob(id, currentUser)

		//we need an struct to return multiple files to the browser with
		var snapData []map[string]string = make([]map[string]string, 0)
		var mediaType = ""
		if client.IsImage(data) {
			mediaType = "image"
		} else if client.IsVideo(data) {
			mediaType = "video"
		} else {
			//in some semi-rare cases, data is returned zipped
			zippedData := client.Unzip(data)
			for _, zipDat := range zippedData {
				if client.IsOverlay(zipDat) {
					snapData = append(snapData, map[string]string{"type": "image", "data": base64.StdEncoding.EncodeToString(zipDat)})
				} else {
					snapData = append(snapData, map[string]string{"type": "video", "data": base64.StdEncoding.EncodeToString(zipDat)})
				}
			}
			r.JSON(200, snapData)
			return
		}
		snapData = append(snapData, map[string]string{"type": mediaType, "data": base64.StdEncoding.EncodeToString(data)})
		r.JSON(200, snapData)
		return
	})

	m.Post("/api/story", authorize, func(res http.ResponseWriter, req *http.Request, currentUser client.User, r render.Render) {
		//find the user in the friend array
		var stories []client.FriendStory
		for _, friend := range currentUser.FriendStories.Friends {
			if friend.Username == req.FormValue("name") {
				stories = friend.FriendStories
			}
		}

		if len(stories) == 0 {
			res.WriteHeader(404) //user not found
			return
		}

		//we need an struct to return multiple files to the browser with
		var storyData []map[string]string = make([]map[string]string, 0)

		for _, story := range stories {
			data := client.GetStory(
				story.Stories.MediaId,
				currentUser.AuthToken,
				string(client.DecodeBase64(story.Stories.MediaKey)),
				string(client.DecodeBase64(story.Stories.MediaIv)))

			if client.IsImage(data) {
				storyData = append(storyData, map[string]string{"type": "image", "data": base64.StdEncoding.EncodeToString(data)})
			} else if client.IsVideo(data) {
				storyData = append(storyData, map[string]string{"type": "video", "data": base64.StdEncoding.EncodeToString(data)})
			} else {
				//in some semi-rare cases, data is returned zipped
				zippedData := client.Unzip(data)
				for _, zipDat := range zippedData {
					fmt.Println("Data=", zipDat[0:2])
					if client.IsOverlay(zipDat) {
						fmt.Println("Is Overlay")
						storyData = append(storyData, map[string]string{"type": "image", "data": base64.StdEncoding.EncodeToString(zipDat)})
					} else {
						storyData = append(storyData, map[string]string{"type": "video", "data": base64.StdEncoding.EncodeToString(zipDat)})
					}
				}
			}

		}

		r.JSON(200, storyData)
	})
	m.Run()
}

//The authorize middleware will search the session for a user_id
//if it doesnt find it, it will return BAD_REQUEST
func authorize(w http.ResponseWriter, r *http.Request, session sessions.Session, c martini.Context) {
	userId := session.Get("user_id")

	if userId != nil {
		sUserId := userId.(string)
		if currentUser, ok := users[sUserId]; ok {
			c.Map(currentUser)
			return
		}
	}
	w.WriteHeader(400)

}
