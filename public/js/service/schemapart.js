
import { get as api_get, create as api_create } from '../api/schemapart.js';

const cache = {};

function getKey(schema_id, block_x, block_y, block_z){
  return `${schema_id}/${block_x}/${block_y}/${block_z}`;
}

export function get(schema_id, block_x, block_y, block_z){
  const key = getKey(schema_id, block_x, block_y, block_z);
  if (cache[key]){
    return Promise.resolve(cache[key]);
  } else {
    return api_get(schema_id, block_x, block_y, block_z)
    .then(schemapart => {
      cache[key] = schemapart;
      return schemapart;
    });
  }
}

export function create(schemapart){
  const key = getKey(schemapart.schema_id, schemapart.offset_x, schemapart.offset_y, schemapart.offset_z);
  cache[key] = schemapart;
  return api_create(schemapart);
}
