export function getNodePos(schemapart, x,y,z){
  return z +
		(y * schemapart.data.size.z) +
		(x * schemapart.data.size.y * schemapart.data.size.z);
}


export function isNodeHidden(schemapart,x,y,z){
  if (x<=1 || x>=14 || y<=1 || y>=14 || z<=1 || z>=14){
    // not sure, may be visible
    return false;
  }

  function isTransparent(contentId){
    var nodeName = schemapart.data.node_mapping_rev[contentId];
    return nodeName == "air" || nodeName == "default:water_source";
  }

  if (isTransparent(schemapart.data.node_ids[getNodePos(schemapart,x-1,y,z)]))
    return false;
  if (isTransparent(schemapart.data.node_ids[getNodePos(schemapart,x,y-1,z)]))
    return false;
  if (isTransparent(schemapart.data.node_ids[getNodePos(schemapart,x,y,z-1)]))
    return false;
  if (isTransparent(schemapart.data.node_ids[getNodePos(schemapart,x+1,y,z)]))
    return false;
  if (isTransparent(schemapart.data.node_ids[getNodePos(schemapart,x,y+1,z)]))
    return false;
  if (isTransparent(schemapart.data.node_ids[getNodePos(schemapart,x,y,z+1)]))
    return false;

  return true;
}
