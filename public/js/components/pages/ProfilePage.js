import UserProfile from '../user/UserProfile.js';
import CreateToken from '../user/CreateToken.js';

export default {
	components: {
		"user-profile": UserProfile,
		"create-token": CreateToken
	},
	template: /*html*/`
		<div>
			<user-profile/>
			<br>
			<create-token/>
		</div>
	`
};
