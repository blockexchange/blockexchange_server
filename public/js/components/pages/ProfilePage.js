import UserProfile from '../user/UserProfile.js';
import AccessToken from '../user/AccessToken.js';

export default {
	components: {
		"user-profile": UserProfile,
		"access-token": AccessToken
	},
	template: /*html*/`
		<div>
			<user-profile/>
			<br>
			<access-token/>
		</div>
	`
};
