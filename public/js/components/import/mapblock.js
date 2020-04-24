
export function createMapblock(offset_x, offset_y, offset_z){
  const mapblock = {
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

  //fill holes
  for (let i=0; i<4096; i++){
    mapblock.data.node_ids[i] = 0;
    mapblock.data.param1[i] = 0;
    mapblock.data.param2[i] = 0;
  }

  return mapblock;
}
