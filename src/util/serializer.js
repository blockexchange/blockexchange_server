const zlib = require('zlib');


/*
Input:

data = {
  node_ids: [0,0,0,...],
  param1: [15,15,15,...],
  param2: [0,0,0...],
  node_mapping: {
    "air": 0,
    "default:dirt": 200
  },
  metadata: {
    "(0,0,0)": {
      fields: {},
      inventories: {}
    }
  },
  size = {
    x: 16,
    y: 16,
    z: 16
  }
}

Output:

{
  data: Buffer(),
  metadata: Buffer()
}
*/
module.exports.serialize = function(data) {
  const data_buffer = Buffer.alloc(
    (data.node_ids.length * 2) +
    data.param1.length +
    data.param2.length
  );

  let offset = 0;
  for (let i=0; i<data.node_ids.length; i++){
    data_buffer.writeInt16LE(data.node_ids[i], offset);
    offset += 2;
  }

  for (let i=0; i<data.param1.length; i++){
    data_buffer.writeInt8(data.param1[i], offset);
    offset++;
  }

  for (let i=0; i<data.param2.length; i++){
    data_buffer.writeInt8(data.param2[i], offset);
    offset++;
  }

  const metadata_json = JSON.stringify({
    node_mapping: data.node_mapping,
    metadata: data.metadata,
    size: data.size
  });

  return {
    metadata: zlib.gzipSync(metadata_json),
    data: zlib.gzipSync(data_buffer)
  };
};
