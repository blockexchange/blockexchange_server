import { setup as setupControls, update as updateControls } from './controls.js';
import { init as initColormapping } from './colormapping.js';
import drawMapblock from './render.js';
import iterator from './iterator.js';

const height = 640;
const width = 800;

export default {
  view: function(vnode){
    console.log("oviewninit", vnode);
    return m("div");
  },
  oncreate: function(vnode) {
    console.log("oncreate", vnode);

		const schema = vnode.attrs.schema;
    const camera = new THREE.PerspectiveCamera( 45, height / width, 2, 2000 );
    camera.position.z = -150;
    camera.position.x = -150;
    camera.position.y = 100;

    const scene = new THREE.Scene();
		vnode.state.scene = scene;
		vnode.state.active = true;

    const renderer = new THREE.WebGLRenderer({
      antialias: false,
      precision: "lowp"
    });

		vnode.state.renderer = renderer;

    renderer.setPixelRatio( window.devicePixelRatio );
    renderer.setSize( width, height );
    document.body.appendChild( renderer.domElement );

    setupControls(camera, renderer, render);

		const it = iterator(schema);

		function fetchNextMapblock(){
			if (!vnode.state.active){
				return;
			}
			const pos = it();

			if (pos){
				drawMapblock(scene, schema, pos.x, pos.y, pos.z)
				.then(render)
				.then(() => setTimeout(fetchNextMapblock, 500));
			}
		}

		initColormapping().then(fetchNextMapblock);

    function render(){
    	renderer.render( scene, camera );
    }

		animate();

    function animate() {
			if (!vnode.state.active){
				return;
			}
    	requestAnimationFrame( animate );
    	updateControls();
    }
  },
  onbeforeupdate: function() {
    return false;
  },
  onremove: function(vnode) {
    console.log("onremove", vnode);
		vnode.state.renderer.dispose();
		vnode.state.scene.dispose();
		vnode.state.active = false;
  }
};
