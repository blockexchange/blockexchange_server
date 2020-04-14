import { getColor } from './colormapping.js';

var materialCache = {};

export function getMaterial(nodeName){
  var material = materialCache[nodeName];
  if (!material) {
    var colorObj = getColor(nodeName);

    if (!colorObj){
      return;
    }

		if (nodeName == "default:stone") {
			let texture = new THREE.TextureLoader().load("textures/default_stone.png");
			material = new THREE.MeshBasicMaterial( { map: texture } );

		}	else if (nodeName == "default:dirt") {
			let texture = new THREE.TextureLoader().load("textures/default_dirt.png");
			material = new THREE.MeshBasicMaterial( { map: texture } );

		}	else if (nodeName == "default:dirt_with_grass") {
			let texture = new THREE.TextureLoader().load("textures/default_grass.png");
			material = new THREE.MeshBasicMaterial( { map: texture } );

		}	else if (nodeName == "default:sand") {
			let texture = new THREE.TextureLoader().load("textures/default_sand.png");
			material = new THREE.MeshBasicMaterial( { map: texture } );

		} else {
			let color = new THREE.Color( colorObj.r/256, colorObj.g/256, colorObj.b/256 );
			material = new THREE.MeshBasicMaterial( { color: color } );

		}

		if (nodeName == "default:water_source"){
			material.transparent = true;
			material.opacity = 0.5;
		}

    materialCache[nodeName] = material;
  }

  return material;
}
