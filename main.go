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

	"github.com/adriangonzy/websocket-balls/game"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveCanvas(w http.ResponseWriter, r *http.Request) {
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

func initSimulation() {
	fmt.Println("INIT SIMULATION")
	balls, e := json.Marshal(game.Balls)
	if e != nil {
		fmt.Errorf("%v", e)
	}

	fmt.Printf("%s\n", balls)

	h.broadcast <- balls
}

func startSimulation() {
	fmt.Println("START SIMULATION")
	stream := game.Start()

	go func() {
		for balls := range stream {
			bytes, e := json.Marshal(balls)
			if e != nil {
				fmt.Errorf("%v", e)
			}
			h.broadcast <- bytes
		}
	}()
}

func stopSimulation() {
	fmt.Println("STOP SIMULATION")
}

const (
	create = "create"
	start  = "start"
	stop   = "stop"
)

type CMD struct {
	Action string `json:"action"`
}

func gameCallBack(message []byte) {
	var cmd CMD

	fmt.Printf("MSG RECEIVED: %s", message)
	e := json.Unmarshal(message, &cmd)
	if e != nil {
		fmt.Errorf("%s", e)
		return
	}

	fmt.Printf("CMD PARSED: %s", cmd)

	switch cmd.Action {
	case start:
		startSimulation()
	case create:
		initSimulation()
	case stop:
		stopSimulation()
	}
}

// serverWs handles webocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, callback: gameCallBack}

	h.register <- c
	go c.writePump()
	c.readPump()
}

func main() {
	flag.Parse()

	go h.run()

	http.HandleFunc("/balls", serveCanvas)
	http.HandleFunc("/ws", serveWs)

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
