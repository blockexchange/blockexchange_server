import { getMaterial } from './material.js';
import { getMapblock } from './api.js';
import { isNodeHidden, getNodePos } from './utils.js';


export function drawMapblock(scene, posx, posy, posz){
  return getMapblock(posx, posy, posz)
  .then(function(mapblock){
    if (!mapblock)
      return;

    if (mapblock.blockmapping.length == 1 && mapblock.blockmapping[0] == "air"){
      return;
    }

    var nodenameGeometriesMap = {}; // nodeName => [matrix, matrix, ...]

		for (var x=0; x<16; x++){
			for (var y=0; y<16; y++){
	      for (var z=0; z<16; z++){
          if (isNodeHidden(mapblock, x,y,z)){
            //skip hidden node
            continue;
          }

					var i = getNodePos(x,y,z);
          var contentId = mapblock.contentid[i];
          var nodeName = mapblock.blockmapping[contentId];

          var matrix = new THREE.Matrix4()
            .makeTranslation(
              x + (posx*16),
              y + (posy*16),
              z + (posz*16)
            );

          var list = nodenameGeometriesMap[nodeName];
          if (!list){
            list = [];
            nodenameGeometriesMap[nodeName] = list;
          }

          list.push(matrix);
        }
      }
    }

    Object.keys(nodenameGeometriesMap).forEach(function(nodeName){
      var material = getMaterial(nodeName);

      if (material){
				var list = nodenameGeometriesMap[nodeName];
				var geometry = new THREE.BoxGeometry(1,1,1);
				var mesh = new THREE.InstancedMesh( geometry, material, list.length );
				for (var i=0; i<list.length; i++){
					mesh.setMatrixAt(i, list[i]);
				}

        scene.add( mesh );
      }
    });
  });
}
