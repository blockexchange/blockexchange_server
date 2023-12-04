export default {
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Blockexchange</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0">
					<li class="nav-item">
						<router-link to="/login" class="nav-link">
							<i class="fa fa-sign-in"></i> Login
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/profile" class="nav-link">
							<i class="fa fa-user"></i> Profile
						</router-link>
					</li>
				</ul>
			</div>
		</nav>
	`
};
