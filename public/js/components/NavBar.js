import LoginStatus from './login/LoginStatus.js';

export default {
	components: {
		"login-status": LoginStatus
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<router-link to="/" class="navbar-brand">Block exchange</router-link>
			<ul class="navbar-nav mr-auto">
				<li class="nav-item">
					<router-link to="/" class="nav-link">
						<i class="fa fa-question"></i> About
					</router-link>
				</li>
				<li class="nav-item">
					<router-link to="/login" class="nav-link">
						<i class="fa fa-sign-in"></i> Login
					</router-link>
				</li>
				<li class="nav-item">
					<router-link to="/mod" class="nav-link">
						<i class="fa fa-download"></i> Mod/Installation
					</router-link>
				</li>
				<li class="nav-item">
					<router-link to="/users" class="nav-link">
						<i class="fa fa-users"></i> Users
					</router-link>
				</li>
				<li class="nav-item">
					<router-link to="/search" class="nav-link">
						<i class="fa fa-search"></i> Search
					</router-link>
				</li>
			</ul>
			<form class="form-inline my-2 my-lg-0">
				<login-status/>
			</form>
	</nav>
	`
};
