
let controls;

export function setup(camera, renderer, changecallback){
  controls = new THREE.TrackballControls( camera, renderer.domElement );
  controls.rotateSpeed = 2.0;
  controls.zoomSpeed = 1.2;
  controls.panSpeed = 0.8;

  controls.noZoom = false;
  controls.noPan = false;

  controls.staticMoving = true;
  controls.dynamicDampingFactor = 0.3;

  controls.keys = [ 65, 83, 68 ];

  controls.addEventListener('change', changecallback);
}

export function update(){
  controls.update();
}
