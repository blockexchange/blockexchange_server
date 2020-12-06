import LocalLogin from '../login/LocalLogin.js';

export default {
	components: {
		'local-login': LocalLogin
	},
	template: /*html*/`
		<div class="row">
			<div class="col-md-8">
				<div class="card">
				  <div class="card-header">
				    Local login
				  </div>
				  <div class="card-body">
				    <h5 class="card-title">Login</h5>
						<local-login/>
				  </div>
				</div>
			</div>
			<div class="col-md-4">
				<div class="card">
					<div class="card-header">
						Sign up
					</div>
					<div class="card-body">
						<router-link to="/register" class="btn btn-secondary">
							Register
						</router-link>
					</div>
				</div>
				<br>
				<div class="card">
					<div class="card-header">
						External login
					</div>
					<div class="card-body">
						<a href="https://github.com/login/oauth/authorize?client_id=68c2728e22f3a4b02dc0" class="btn btn-secondary">
							<i class="fa fa-github"></i>
							Login with Github
						</a>
					</div>
				</div>
			</div>
		</div>
	`
};
