# Go Snap!

A Pure Go library for the Snapchat API

The implementation is partially based on [pysnap](https://github.com/martinp/pysnap) by martinp and the [full disclosure](http://gibsonsec.org/snapchat/fulldisclosure/) by GibSec.


### Features

* [x] Login/Logout and list out your recent snaps
* [x] Retreive both picture and video snaps, in your browser
* [x] Download the pictures or snaps, from the browser
* [x] Store the user in the session for easy reterival
* [x] List & Decrypt Story images and videos
* [ ] Send snaps
* [ ] Store users (and snaps?) in a database


 ![GoSnap](github.com/jamieomatthews/gosnap/public/img/gosnap.png)

### Installation

```bash
git clone https://github.com/jamieomatthews/gosnap.git
go run app.go
```

### Components

The Gosnap client package has zero external dependencies, and can function as is.  I wanted Gosnap to be easier to use than most of the command line utilties out there, so I wrote a small webapp in Martini that lets you browse your snaps, and view them.

### Saving Snaps

Currently, the only client is a web based client.  It would be trivial to write a client that simply saved the snaps to disk, but for now, if you want to save a snap, you can simply right click on the image or video, and save.

### Contributing

Contributions are welcome, I would definitely like to finish out the feature set, and improve the user interface where applicable.  Check the feature list if you are are looking for what to do next, or if you have an idea of your own.  To contribute, just send a PR!