import { getMaterial } from './material.js';
import { isNodeHidden, getNodePos } from './utils.js';


export default function(scene, schema, posx, posy, posz){

	const block_x = Math.floor(posx * schema.part_length);
	const block_y = Math.floor(posy * schema.part_length);
	const block_z = Math.floor(posz * schema.part_length);

  return m.request(`api/schemapart/${schema.id}/${block_x}/${block_y}/${block_z}`)
  .then(function(mapblock){
    if (!mapblock)
      return;

		if (Object.keys(mapblock.data.node_mapping).length == 1 && mapblock.data.node_mapping.air) {
      return;
    }

		// create reverse mapping: nodeId -> nodeName
		const node_mapping_rev = {};
		Object.keys(mapblock.data.node_mapping).forEach(nodeName => {
			const nodeId = mapblock.data.node_mapping[nodeName];
			node_mapping_rev[nodeId] = nodeName;
		});
		mapblock.data.node_mapping_rev = node_mapping_rev;

    var nodenameGeometriesMap = {}; // nodeName => [matrix, matrix, ...]

		for (var x=0; x<16; x++){
			for (var y=0; y<16; y++){
	      for (var z=0; z<16; z++){
          if (isNodeHidden(mapblock, x,y,z)){
            //skip hidden node
            continue;
          }

					var i = getNodePos(x,y,z);
          var contentId = mapblock.data.node_ids[i];
          var nodeName = mapblock.data.node_mapping_rev[contentId];

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
