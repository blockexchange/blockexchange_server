import { setup as setupControls, update as updateControls } from './controls.js';

export default {
  view: function(vnode){
    console.log("oviewninit", vnode);
    return m("div");
  },
  oncreate: function(vnode) {
    console.log("oncreate", vnode);
    const camera = new THREE.PerspectiveCamera( 45, window.innerWidth / window.innerHeight, 2, 2000 );
    camera.position.z = -150;
    camera.position.x = -150;
    camera.position.y = 100;

    const scene = new THREE.Scene();

    const renderer = new THREE.WebGLRenderer({
      antialias: false,
      precision: "lowp"
    });

    renderer.setPixelRatio( window.devicePixelRatio );
    renderer.setSize( window.innerWidth, window.innerHeight );
    document.body.appendChild( renderer.domElement );

    setupControls(camera, renderer, render);

    function render(){
    	renderer.render( scene, camera );
    }

    function animate() {
    	requestAnimationFrame( animate );
    	updateControls();
    }
  },
  onbeforeupdate: function() {
    return false;
  },
  onremove: function(vnode) {
    console.log("onremove", vnode);
  }
};
