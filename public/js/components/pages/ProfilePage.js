import UserProfile from '../user/UserProfile.js';
import AccessToken from '../user/AccessToken.js';
import UserSchemas from '../search/UserSchemas.js';
import store from '../../store/login.js';

export default {
	data: () => store,
	components: {
		"user-profile": UserProfile,
		"access-token": AccessToken,
		"user-schemas": UserSchemas
	},
	template: /*html*/`
		<div>
			<user-profile/>
			<br>
			<access-token/>
			<br>
			<user-schemas :user_name="claims.username"/>
		</div>
	`
};
