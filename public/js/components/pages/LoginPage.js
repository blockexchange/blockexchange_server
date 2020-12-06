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
						<a class="btn btn-secondary">Register</a>
					</div>
				</div>
				<br>
				<div class="card">
					<div class="card-header">
						External login
					</div>
					<div class="card-body">
						<a class="btn btn-secondary">Login with Github</a>
					</div>
				</div>
			</div>
		</div>
	`
};
