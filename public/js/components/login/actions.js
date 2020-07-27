import { set_token } from '../../store/token.js';
import { request_token } from '../../api/token.js';

import state from './state.js';

export function login(){
	state.message = null;
	request_token(state.username, state.password)
	.then(token => set_token(token))
	.catch(e => state.message = e.message);
}

export function temp_login(){
	state.message = null;
	request_token("temp", "temp")
	.then(token => set_token(token))
	.catch(e => state.message = e.message);
}

export function logout(){
	set_token(null);
}
