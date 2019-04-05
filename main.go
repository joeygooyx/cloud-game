package main

import (
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"net/http"

	// "time"

	"github.com/giongto35/game-online/ui"
	"github.com/giongto35/game-online/util"
	"github.com/giongto35/game-online/webrtc"
	"github.com/gorilla/mux"
)

// var webRTC *webrtc.WebRTC
var width = 256
var height = 240
var gameName = "Contra.nes"

// var FPS = 60

func init() {
}

func startGame(path string, imageChannel chan *image.RGBA, inputChannel chan int, webRTC *webrtc.WebRTC) {
	ui.Run([]string{path}, imageChannel, inputChannel, webRTC)
}

func main() {
	fmt.Println("http://localhost:8000")

	router := mux.NewRouter()
	router.HandleFunc("/", getWeb).Methods("GET")
	router.HandleFunc("/session", postSession).Methods("POST")

	http.ListenAndServe(":8000", router)
}

func getWeb(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadFile("./index.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(bs)
}

func postSession(w http.ResponseWriter, r *http.Request) {
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()

	webRTC := webrtc.NewWebRTC()

	localSession, err := webRTC.StartClient(string(bs), width, height)
	if err != nil {
		log.Fatalln(err)
	}

	imageChannel := make(chan *image.RGBA, 100)
	go screenshotLoop(imageChannel, webRTC)
	go startGame("games/"+gameName, imageChannel, webRTC.InputChannel, webRTC)

	w.Write([]byte(localSession))
}

// func screenshotLoop(imageChannel chan *image.RGBA) {
func screenshotLoop(imageChannel chan *image.RGBA, webRTC *webrtc.WebRTC) {
	for image := range imageChannel {
		// Client stopped
		if webRTC.IsClosed() {
			break
		}

		// encode frame
		if webRTC.IsConnected() {
			yuv := util.RgbaToYuv(image)
			webRTC.ImageChannel <- yuv
		}
	}
}