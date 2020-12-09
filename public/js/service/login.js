import { request_token } from '../api/token.js';
import loginstore from '../store/login.js';

const STORAGE_KEY = "blockexchange";

export default {
	login(username, password){
		return request_token(username, password)
		.then(t => {
			loginstore.token = t;
			loginstore.loggedIn = true;
			loginstore.username = username;

			this.persist();
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
