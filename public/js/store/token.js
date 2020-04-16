
var token = localStorage.blockexchange_token;
var claims = null;

function parse_claims(){
  if (token){
    claims = JSON.parse(atob(token.split(".")[1]));
  } else {
    claims = null;
  }
}

parse_claims();

export function get_token(){
  return token;
}

export function get_claims(){
  return claims;
}

export function set_token(t){
  if (!t){
    token = null;
    delete localStorage.blockexchange_token;
  } else {
    token = t;
    localStorage.blockexchange_token = t;
  }
  parse_claims();
}
