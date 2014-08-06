package main

import (
	"fmt"
	"log"
	"strconv"

	"encoding/json"
	"net/http"

	"github.com/adriangonzy/websocket-balls/game"
	"github.com/adriangonzy/websocket-balls/ws"
)

var conn *ws.Connection
var sim *game.Simulation

func bindSimulationControls() {
	http.HandleFunc("/simulation/start", startSimulation)
	http.HandleFunc("/simulation/stop", stopSimulation)
	http.HandleFunc("/ws", serveWs)
}

func startSimulation(w http.ResponseWriter, r *http.Request) {

	if conn == nil {
		http.Error(w, "Must upgrade to websocket connection before starting simulation ", http.StatusInternalServerError)
		return
	}

	// parse query param
	ballsCount, err := strconv.Atoi(r.URL.Query().Get("balls_count"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// init simulation with given number of balls
	sim := game.NewSimulation(ballsCount)
	sim.Start()

	go func() {
		for balls := range sim.Balls {
			conn.Send <- serializeBalls(balls)
		}
	}()
}

func serializeBalls(balls []*game.Ball) []byte {
	b, e := json.Marshal(balls)
	if e != nil {
		fmt.Errorf("%v", e)
		return nil
	}
	return b
}

func stopSimulation(w http.ResponseWriter, r *http.Request) {
	if sim == nil {
		http.Error(w, "Must start simulation before stopping it", http.StatusInternalServerError)
		return
	}

	sim.Stop()
}

// serverWs handles webocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	websocket, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error Upgrading", err)
		return
	}
	conn = ws.NewConnection(make(chan []byte, 256), websocket)
	conn.Start()
	log.Println("Connection STARTED")
}