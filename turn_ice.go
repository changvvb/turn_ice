package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type turnServer struct {
	Username   string   `json:"username"`
	Credential string   `json:"credential"`
	Urls       []string `json:"urls"`
}

type stunServer struct {
	Urls []string `json:"urls"`
}

type iceServer struct {
	Servers []interface{} `json:"iceServers"`
}

type servers struct {
	TurnServer turnServer
	StunServer stunServer
}

func ice(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE")

	r.ParseForm()
	username := template.HTMLEscapeString(r.Form.Get("username"))
	credential := template.HTMLEscapeString(r.Form.Get("key"))
	log.Println("ICE: username:", username, " key:", credential)
	/*
		enc := json.NewEncoder(w)
		if err := enc.Encode(iceserver); err != nil {
			log.Println(err)
		}
	*/
	b, err := json.Marshal(iceserver)
	fmt.Fprint(w, string(b))
	if err != nil {
		log.Println(err)
	}
}

func turn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := template.HTMLEscapeString(r.Form.Get("username"))
	credential := template.HTMLEscapeString(r.Form.Get("key"))
	log.Println("TURN: username:", username, " key:", credential)
	enc := json.NewEncoder(w)
	if err := enc.Encode(turnserver); err != nil {
		log.Println(err)
	}
}

var iceserver iceServer
var turnserver turnServer
var stunserver stunServer

func main() {
	fmt.Println("HEHEHEHHEHE")

	turnserver.Urls = []string{"turn:115.29.55.106:3478?transport=udp", "turn:115.29.55.106:3478?transport=tcp", "turn:115.29.55.106:3479?transport=udp", "turn:115.29.55.106:3479?transport=tcp"}
	stunserver.Urls = []string{"stun:stun.115.29.55.106"}
	iceserver.Servers = append(iceserver.Servers, &turnserver)
	iceserver.Servers = append(iceserver.Servers, &stunserver)
	turnserver.Username = "changvvb"
	turnserver.Credential = "changvvb"

	http.HandleFunc("/ice", ice)
	http.HandleFunc("/turn", turn)
	http.ListenAndServe(":9000", nil)
}
