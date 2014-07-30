// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
	"text/template"

	"github.com/adriangonzy/websocket-balls/game"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("front/chat.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

func servePlayground(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "playground.html")
}

func compressBalls(balls []*game.Ball) []interface{} {
	bs := make([]interface{}, len(balls)*2)
	for i, j := 0, 0; i < len(balls) && j < len(balls)*2; i++ {
		bs[j] = balls[i].Position.X
		bs[j+1] = balls[i].Position.Y
		j += 2
	}
	return bs
}

func startSimulation(w http.ResponseWriter, r *http.Request) {
	fmt.Println("start simulation")
	balls, e := json.Marshal(game.Balls)
	if e != nil {
		fmt.Errorf("%v", e)
	}

	fmt.Printf("%s\n", balls)

	h.broadcast <- balls

	stream := game.Start()

	go func() {
		for balls := range stream {
			bytes, e := json.Marshal(compressBalls(balls))
			if e != nil {
				fmt.Errorf("%v", e)
			}
			h.broadcast <- bytes
		}
	}()
}

// serverWs handles webocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	fmt.Println("socket connection established")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
	c.readPump()
}

func main() {
	flag.Parse()
	go h.run()
	//*
	http.HandleFunc("/chat", serveHome)
	http.HandleFunc("/balls", servePlayground)
	//*/
	http.Handle("/front/js/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/ws", serveWs)

	http.HandleFunc("/start", startSimulation)

	/*
		for balls := range stream {
			fmt.Println("===============================================================")
			j, e := json.Marshal(compressBalls(balls))
			if e != nil {
				fmt.Errorf("%v", e)
			}
			fmt.Printf("%s\n", j)
		}
	*/

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
