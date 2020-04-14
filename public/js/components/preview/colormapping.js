let mapping;

export function init(){
  return m.request("colormapping.json")
    .then(_mapping => mapping = _mapping);
}

export function getColor(nodeName){
  return mapping[nodeName];
}
