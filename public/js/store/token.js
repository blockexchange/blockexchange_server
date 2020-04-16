
var token = localStorage.blockexchange_token;

export function get_token(){
  return token;
}

export function get_claims(){
  // TODO
}

export function set_token(t){
  token = t;
  localStorage.blockexchange_token = t;
}
