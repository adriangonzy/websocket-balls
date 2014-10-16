// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"

	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveCanvas(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "balls.html")
}

func serveSlider(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "style/css; charset=utf-8")
	http.ServeFile(w, r, "slider.css")
}

func main() {
	flag.Parse()

	http.HandleFunc("/balls", serveCanvas)
	http.HandleFunc("/slider.css", serveSlider)
	bindSimulationControls()

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
