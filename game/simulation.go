package game

import (
	"fmt"
)

func init() {
	fmt.Println("Simulation initialized")
}

const (
	canvasHeight = 600
	canvasWidth  = 900
	maxVelocity  = 10
	maxRadius    = 50
	maxMass      = 1000
)

func run() (chan []ball) {
	
	for {
		update()
	}
}

  function mainLoop() {
        thisTime = Date.now();
        deltaTime = thisTime - lastTime;

        renderer.draw(context, ballArray);
        simulation.update(deltaTime, ballArray);

        lastTime = thisTime;

        setTimeout(mainLoop, frameTimer);
    }