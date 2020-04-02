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
    data_buffer.writeInt8(data.param1[i] - 128, offset);
    offset++;
  }

  for (let i=0; i<data.param2.length; i++){
    data_buffer.writeInt8(data.param2[i] - 128, offset);
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


/*
Input: {
 metadata: Buffer(),
 data: Buffer()
}

Output: {
 .. same as above input ..
}

*/
module.exports.deserialize = function(data) {
  const metadata_json = zlib.gunzipSync(data.metadata);
  const result = JSON.parse(metadata_json);
  const param_length = result.size.x * result.size.y * result.size.z;
  const buffer_length = (param_length * 2) + param_length + param_length;

  const data_buffer = zlib.gunzipSync(data.data);
  if (buffer_length != data_buffer.length) {
    throw new Error("Unexpected size: " + data_buffer.length + " should be: " + buffer_length + " metadata: " + metadata_json);
  }

  result.node_ids = [];
  let offset = 0;
  for (let i=0; i<param_length; i++){
    result.node_ids.push(data_buffer.readInt16LE(offset));
    offset += 2;
  }

  result.param1 = [];
  for (let i=0; i<param_length; i++){
    result.param1.push(data_buffer.readInt8(offset) + 128);
    offset++;
  }

  result.param2 = [];
  for (let i=0; i<param_length; i++){
    result.param2.push(data_buffer.readInt8(offset) + 128);
    offset++;
  }

  return result;
};
