package main

import (
	"fmt"
	"log"

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

func getConfig(req *http.Request) *game.Config {
	decoder := json.NewDecoder(req.Body)
	fmt.Printf("request body %#v \n", req.Body)
	var t game.Config
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	return &t
}

func startSimulation(w http.ResponseWriter, r *http.Request) {

	if conn == nil {
		http.Error(w, "Must upgrade to websocket connection before starting simulation", http.StatusInternalServerError)
		log.Fatal("Must upgrade to websocket connection before starting simulation")
		return
	}

	c := getConfig(r)

	// init simulation with given number of balls
	sim = game.NewSimulation(c)
	sim.Start()

	go func() {
		for balls := range sim.Emit {
			conn.Send <- serializeBalls(balls)
		}
	}()
}

func serializeBalls(balls [][]interface{}) []byte {
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
