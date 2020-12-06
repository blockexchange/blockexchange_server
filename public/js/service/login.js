import { request_token } from '../api/token.js';
import store from '../store/token.js';

export default {
	login(username, password){
		return request_token(username, password)
		.then(t => {
			store.token = t;
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
		store.token = null;
	}
};
