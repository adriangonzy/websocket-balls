/*###########################################################################
  # Code created by Adam Brookes for adambrookesprojects.co.uk - 06/10/2013 #
  #  It is unencumbered by copyrights and patents and we can use it freely, #
  # but we can only assert our own Intellectual Property rights on derived  #
  #  works: the original work remains free for public use                   #
  ########################################################################### 

 #######################################################
 # Class used to manage Canvas Renderer and Simulation #
 #######################################################*/
if (window.addEventListener) window.addEventListener('load', onLoad, false);

function onLoad() {
    var canvas;
    var canvasHeigth = 600;
    var canvasWidth = 900;
    var context;
    
    var renderer = new Renderer('#CAE1FF'); // takes colour for canvas.
    var simulation;
    var ballArray = new Array();
    var numberOfBalls = 40;

    // frameRate Variables.
    var frameRate = 60;
    var frameTimer = 1000 / frameRate;

    // DeltaTime variables.
    var lastTime = Date.now(); // inistalise lastTime.
    var thisTime;
    var deltaTime;
    
    function initialiseCanvas() {
        //find the canvas element using its id attribute.
        canvas = document.getElementById('canvas');
        //once canvas is created, create the simulation passing the width and height of canvas
        canvas.width = canvasWidth;
        canvas.height = canvasHeigth;
        simulation = new Simulation(canvas.width, canvas.height);

        /*########## Error checking to see if canvas is supported ############## */
        if (!canvas) {
            alert('Error: cannot find the canvas element!');
            return;
        }
        if (!canvas.getContext) {
            alert('Error: no canvas.getContent!');
            return;
        }
        context = canvas.getContext('2d');
        if (!context) {
            alert('Error: failed to getContent');
            return;
        }
        createBalls();
        mainLoop(); // enter the main loop.
    }

    function random(num) {
        return Math.floor(Math.random() * num) + 1;
    }

    function randomColor() {
        return Math.floor(Math.random() * parseInt("ffffff", 16)).toString(16);
    }

    function createBalls() {
        /* Ball takes X | Y | radius | Mass| vX | vY | colour */
        for (var i = 0;i < numberOfBalls; i++) {
            ballArray.push(new ball(random(canvasWidth), 
                                    random(canvasHeigth),
                                    random(3),
                                    random(1000),
                                    random(5), 
                                    random(5), 
                                    randomColor()));
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

    initialiseCanvas();
}