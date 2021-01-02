import { request_token } from '../api/token.js';
import loginstore from '../store/login.js';

const STORAGE_KEY = "blockexchange";

export default {
	login(username, password){
		return request_token(username, password)
		.then(t => {
			loginstore.loggedIn = true;
			this.parse_token(t);
			return { success: true };
		})
		.catch(e => {
			return {
				success: false,
				message: e.message
			};
		});
	},

	parse_token(token){
		const payload = JSON.parse(atob(token.split(".")[1]));
		loginstore.token = token;
		loginstore.loggedIn = true;
		loginstore.username = payload.username;
		loginstore.claims = payload;
		this.persist();
	},

	logout(){
		loginstore.token = null;
		loginstore.loggedIn = false;
		loginstore.username = "";
		loginstore.claims = null;
		this.persist();
	},

	// persist login state
	persist(){
			if (loginstore.token){
				// persist to localStorage
				localStorage[STORAGE_KEY] = JSON.stringify(loginstore);
			} else {
				// clear localStorage
				delete localStorage[STORAGE_KEY];
			}
	},

	// restore store from localStorage data
	restoreState(){
		if (!localStorage[STORAGE_KEY]){
			return;
		}

		const data = JSON.parse(localStorage[STORAGE_KEY]);

		if (!data){
			return;
		}

		// TODO: check if token is expired

		Object.keys(data).forEach(key => {
			loginstore[key] = data[key];
		});
	}
};
