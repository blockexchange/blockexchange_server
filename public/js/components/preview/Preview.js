import { init as initColormapping } from './colormapping.js';
import drawMapblock from './render.js';
import iterator from './iterator.js';
import OrbitControls from './orbitcontrols.js';

export default {
  view: function(){
    return m("div", { style: `height: 450px; width: 100%;`});
  },
  oncreate: function(vnode) {
		const schema = vnode.attrs.schema;
    const camera = new THREE.PerspectiveCamera( 70, 2, 1, 1000 );
    camera.position.z = -150;
    camera.position.x = -150;
    camera.position.y = 100;

    const scene = new THREE.Scene();
		scene.background = new THREE.Color();

		vnode.state.scene = scene;
		vnode.state.active = true;

    const renderer = new THREE.WebGLRenderer({
      antialias: false,
      precision: "lowp",

      // canvas to blob: https://stackoverflow.com/questions/12168909/blob-from-dataurl
      preserveDrawingBuffer: true
    });

		vnode.state.renderer = renderer;

    renderer.setPixelRatio( window.devicePixelRatio );
    vnode.dom.appendChild( renderer.domElement );

    const box = renderer.domElement.parentElement.getBoundingClientRect();
    renderer.setSize( box.width, box.height );

		const controls = new OrbitControls( camera, renderer.domElement, renderer.domElement );

		const it = iterator(schema);
		let count = 0;

		function fetchNextMapblock(){
			if (!vnode.state.active){
				return;
			}
			const pos = it();
			count++;

			if (typeof(vnode.attrs.progressCallback) == "function") {
				if (pos)
					vnode.attrs.progressCallback(count / schema.total_parts);
				else
					vnode.attrs.progressCallback(1);
			}

			if (pos){
				drawMapblock(scene, schema, pos.x, pos.y, pos.z)
				.then(() => renderer.render( scene, camera ))
        .then(() => m.redraw())
				.then(() => setTimeout(fetchNextMapblock, 150));
			}
		}

		initColormapping().then(fetchNextMapblock);
		animate();

		controls.addEventListener('change', () => renderer.render( scene, camera ));
		renderer.render( scene, camera );

    function animate() {
			if (!vnode.state.active){
				return;
			}

    	controls.update();
			requestAnimationFrame( animate );
    }
  },
  onbeforeupdate: function() {
    return false;
  },
  onremove: function(vnode) {
		vnode.state.renderer.dispose();
		vnode.state.scene.dispose();
		vnode.state.active = false;
  }
};
