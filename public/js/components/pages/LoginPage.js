import LocalLogin from '../login/LocalLogin.js';

export default {
	components: {
		'local-login': LocalLogin
	},
	template: /*html*/`
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4">
				<local-login/>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
