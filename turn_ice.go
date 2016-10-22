package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type turnServer struct {
	Username   string   `json:"username"`
	Credential string   `json:"password"`
	Urls       []string `json:"uris"`
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
	server := flag.String("server", "127.0.0.1", "the ip of server")
	port := flag.String("port", "9000", "the port to listen")
	flag.Parse()
	IP := *server
	fmt.Println("ip:" + IP + " port:" + *port)
	turnserver.Urls = []string{"turn:" + IP + ":3478?transport=udp", "turn:" + IP + ":3478?transport=tcp", "turn:" + IP + ":3479?transport=udp", "turn:" + IP + ":3479?transport=tcp"}
	stunserver.Urls = []string{"stun:stun." + IP}
	iceserver.Servers = append(iceserver.Servers, &turnserver)
	iceserver.Servers = append(iceserver.Servers, &stunserver)
	turnserver.Username = "testuser"
	turnserver.Credential = "testpass"

	http.HandleFunc("/ice", ice)
	http.HandleFunc("/turn", turn)
	http.ListenAndServe(":"+*port, nil)
}
