
import { get_claims, set_token } from '../../store/token.js';
import { request_token } from '../../api/token.js';

const state = {
  username: "",
  password: "",
  message: null,
  login: function(){
    state.message = null;
    request_token(state.username, state.password)
    .then(token => set_token(token))
    .catch(e => state.message = e.message);
  },
  temp_login: function(){
    state.message = null;
    request_token("temp", "temp")
    .then(token => set_token(token))
    .catch(e => state.message = e.message);
  },
  logout: function(){
    set_token(null);
  },
  isLoggedIn: function(){
    return get_claims();
  }
};

export default state;
