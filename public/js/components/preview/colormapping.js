import colormapping from '../../service/colormapping.js';

let mapping;

export function init(){
  if (!mapping){
    return colormapping()
    .then(_mapping => mapping = _mapping);
  }
  return Promise.resolve();
}

export function getColor(nodeName){
  return mapping[nodeName];
}
