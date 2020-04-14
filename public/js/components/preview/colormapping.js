import { getColormapping } from './api.js';

let mapping;

export function init(){
  return getColormapping()
    .then(_mapping => mapping = _mapping);
}

export function getColor(nodeName){
  return mapping[nodeName];
}
