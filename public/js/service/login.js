import { request_token } from '../api/token.js';
import loginstore from '../store/login.js';

export default {
	login(username, password){
		return request_token(username, password)
		.then(t => {
			loginstore.token = t;
			loginstore.loggedIn = true;
			loginstore.username = username;

			return { success: true };
		})
		.catch(e => {
			return {
				success: false,
				message: e.message
			};
		});
	},

	logout(){
		loginstore.token = null;
		loginstore.loggedIn = false;
		loginstore.username = "";
	}
};
