import store from '../store/token.js';
import { request_token } from '../api/token.js';

export const login = (username, password) => {
  console.log("login", username, password);
  request_token(username, password)
  .then(t => store.token = t)
  .catch(e => console.error(e));
};

export const isLoggedIn = () => {
  console.log("isLoggedIn");
  return !!store.token;
};
