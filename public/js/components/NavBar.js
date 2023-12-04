export default {
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Block exchange</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0">
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
					<li class="nav-item">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> My schemas
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/import" class="nav-link">
							<i class="fa fa-upload"></i> Schema import
						</router-link>
					</li>
					<li class="nav-item">
						<router-link to="/tags" class="nav-link">
							<i class="fa fa-tags"></i> Tags
						</router-link>
					</li>
				</ul>
				<form class="d-flex">
					<div class="btn btn-secondary">
						<router-link to="/profile">
							<i class="fas fa-user"></i>
							<span>
								Logged in as <b>username</b>
							</span>
						</router-link>
					</div>
				</form>
			</div>
		</nav>
	`
};
