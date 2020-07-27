
import { get_claims } from '../../store/token.js';

const state = {
  username: "",
  password: "",
  message: null,
  isLoggedIn: function(){
    return get_claims();
  }
};

export default state;
