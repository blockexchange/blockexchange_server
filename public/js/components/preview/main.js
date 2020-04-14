import { init as initColormapping } from './colormapping.js';
import { drawMapblock } from './render.js';
import { setup as setupControls, update as updateControls } from './controls.js';

var camera, scene, renderer;

init();
animate();

function init() {
	camera = new THREE.PerspectiveCamera( 45, window.innerWidth / window.innerHeight, 2, 2000 );
	camera.position.z = -150;
	camera.position.x = -150;
	camera.position.y = 100;

	scene = new THREE.Scene();

  var min = -7, max = 7;
  var x = min, y = -1, z = min;

  function increment(){
    x++;
    if (x > max){
      z++;
      x = min;
    }
    if (z > max){
      y++;
      z = min;
    }
  }

  var drawLoop = function(){
    if (y >= 3){
      return;
    }

    drawMapblock(scene, x, y, z)
    .then(function(){
      render();
      increment();
      setTimeout(drawLoop, 50);
    });
  };

  initColormapping().then(drawLoop);

	renderer = new THREE.WebGLRenderer({
		antialias: false,
		precision: "lowp"
	});

	renderer.setPixelRatio( window.devicePixelRatio );
	renderer.setSize( window.innerWidth, window.innerHeight );
	document.body.appendChild( renderer.domElement );

  setupControls(camera, renderer, render);

	render();
}

function render(){
	renderer.render( scene, camera );
}

function animate() {
	requestAnimationFrame( animate );
	updateControls();
}
