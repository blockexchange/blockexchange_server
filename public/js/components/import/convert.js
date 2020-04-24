
import { createMapblock } from './mapblock.js';

function pos_to_index(pos){
  return pos.z + (pos.y * 16) + (pos.x * 256);
}

function get_or_create_mapblock(mapblocks_x_slices, parts, mapblock_x, mapblock_y, mapblock_z){
  const offset_x = mapblock_x * 16;
  const offset_y = mapblock_y * 16;
  const offset_z = mapblock_z * 16;

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
    mapblock = createMapblock(offset_x, offset_y, offset_z);

    parts.push(mapblock);
    mapblocks_z_slices[mapblock_z] = mapblock;
  }

  return mapblock;
}

function get_modname(nodename){
  return nodename.split(":")[0];
}

export default function(blocks, stats){
  var max_x = 0, max_y = 0, max_z = 0;

  //mod and node count stats
  stats = stats || {};

  // mapblocks as array
  const parts = [];

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

    const mapblock = get_or_create_mapblock(mapblocks_x_slices, parts, mapblock_x, mapblock_y, mapblock_z);

    //TODO: non-uniform block sizes (other than 16 blocks wide)
    const index = pos_to_index({
      x: block.x - offset_x,
      y: block.y - offset_y,
      z: block.z - offset_z
    });

    if (block.meta && (
        (block.meta.fields && Object.keys(block.meta.fields).length) ||
        (block.meta.inventory && Object.keys(block.meta.inventory).length)
      )){
      // Assign metadata
      const pos_str = `(${block.x},${block.y},${block.z})`;
      mapblock.data.metadata.meta[pos_str] = block.meta;
    }

    mapblock.data.param1[index] = block.param1;
    mapblock.data.param2[index] = block.param2;

    var node_id;
    if (mapblock.data.node_mapping[block.name]){
      // nodename exists
      node_id = mapblock.data.node_mapping[block.name];
    } else {
      // create new entry
      node_id = mapblock.data.max_nodeid++;
      mapblock.data.node_mapping[block.name] = node_id;
    }

    mapblock.data.node_ids[index] = node_id;

    // update stats
    const modname = get_modname(block.name);
    if (stats[modname]){
      stats[modname]++;
    } else {
      stats[modname] = 1;
    }
  });

  // fill in empty mapblocks
  for (let mapblock_x=0; mapblock_x<Math.ceil(max_x / 16); mapblock_x++){
    for (let mapblock_y=0; mapblock_y<Math.ceil(max_y / 16); mapblock_y++){
      for (let mapblock_z=0; mapblock_z<Math.ceil(max_z / 16); mapblock_z++){
        get_or_create_mapblock(mapblocks_x_slices, parts, mapblock_x, mapblock_y, mapblock_z);
      }
    }
  }

  return {
    stats: stats,
    parts: parts,
    max_x: max_x,
    max_y: max_y,
    max_z: max_z
  };
}
