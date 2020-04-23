
function pos_to_index(pos){
  return pos.x + (pos.y * 16) + (pos.z * 256);
}

export default function(blocks){
  var max_x = 0, max_y = 0, max_z = 0;

  // mapblocks as array
  const mapblocks = [];

  // mapblocks as 3 dimensional array indexed by order (0,1,2)
  const mapblocks_x_slices = [];

  blocks.forEach(function(block){
    const mapblock_x = Math.floor(block.x / 16);
    const mapblock_y = Math.floor(block.y / 16);
    const mapblock_z = Math.floor(block.z / 16);

    const offset_x = mapblock_x * 16;
    const offset_y = mapblock_y * 16;
    const offset_z = mapblock_z * 16;

    max_x = Math.max(max_x, offset_x + 15);
    max_y = Math.max(max_y, offset_y + 15);
    max_z = Math.max(max_z, offset_z + 15);

    var mapblocks_y_slices = mapblocks_x_slices[mapblock_x];
    if (!mapblocks_y_slices){
      mapblocks_y_slices = [];
      mapblocks_x_slices[mapblock_x] = mapblocks_y_slices;
    }

    var mapblocks_z_slices = mapblocks_y_slices[mapblock_y];
    if (!mapblocks_z_slices){
      mapblocks_z_slices = [];
      mapblocks_y_slices[mapblock_y] = mapblocks_z_slices;
    }

    var mapblock = mapblocks_z_slices[mapblock_z];
    if (!mapblock){
      mapblock = {
        schema_id: 0, //TODO
        offset_x: offset_x,
        offset_y: offset_y,
        offset_z: offset_z,
        data: {
          node_ids: new Array(4096),
          param1: new Array(4096),
          param2: new Array(4096),
          node_mapping: { air: 0 },
          max_nodeid: 1,
          metadata: {
            meta: {},
            timers: {}
          },
          size: {
            x: 16,
            y: 16,
            z: 16
          }
        },
      };
      mapblocks.push(mapblock);
      mapblocks_z_slices[mapblock_z] = mapblock;
    }

    //TODO: non-uniform block sizes (other than 16 blocks wide)
    const index = pos_to_index(block);

    if (block.meta){
      mapblock.data.metadata.meta = block.meta;
    }

    mapblock.data.param1[index] = block.param1;
    mapblock.data.param2[index] = block.param2;

    var node_id;
    if (mapblock.data.node_ids[block.name]){
      // nodename exists
      node_id = mapblock.data.node_ids[block.name];
    } else {
      // create new entry
      node_id = mapblock.data.max_nodeid++;
      mapblock.data.node_ids[block.name] = node_id;
    }

    mapblock.data.node_ids[index] = node_id;
  });


  return {
    parts: mapblocks,
    max_x: max_x,
    max_y: max_y,
    max_z: max_z
  };
}
