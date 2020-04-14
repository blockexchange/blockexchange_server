export function getNodePos(x,y,z){
  return x + (y * 16) + (z * 256);
}


export function isNodeHidden(mapblock,x,y,z){
  if (x<=1 || x>=14 || y<=1 || y>=14 || z<=1 || z>=14){
    // not sure, may be visible
    return false;
  }

  function isTransparent(contentId){
    var nodeName = mapblock.blockmapping[contentId];
    return nodeName == "air" || nodeName == "default:water_source";
  }

  if (isTransparent(mapblock.contentid[getNodePos(x-1,y,z)]))
    return false;
  if (isTransparent(mapblock.contentid[getNodePos(x,y-1,z)]))
    return false;
  if (isTransparent(mapblock.contentid[getNodePos(x,y,z-1)]))
    return false;
  if (isTransparent(mapblock.contentid[getNodePos(x+1,y,z)]))
    return false;
  if (isTransparent(mapblock.contentid[getNodePos(x,y+1,z)]))
    return false;
  if (isTransparent(mapblock.contentid[getNodePos(x,y,z+1)]))
    return false;

  return true;
}
