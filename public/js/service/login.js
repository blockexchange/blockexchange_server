import store from '../store/token.js';
import { request_token } from '../api/token.js';

export const login = (username, password) => {
  request_token(username, password)
  .then(t => store.setToken(t))
  .catch(e => console.error(e));
};
