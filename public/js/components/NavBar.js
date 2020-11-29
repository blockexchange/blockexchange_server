export default {
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<router-link to="/" class="navbar-brand">Block exchange</router-link>
			<div class="navbar-collapse collapse">
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
						<a class="nav-link" href="#!/mod">
							<i class="fa fa-download"></i> Mod/Installation
						</a>
					</li>
					<li class="nav-item">
						<a class="nav-link" href="#!/users">
							<i class="fa fa-users"></i> Users
						</a>
					</li>
					<li class="nav-item">
						<router-link to="/search" class="nav-link">
							<i class="fa fa-search"></i> Search
						</router-link>
					</li>
				</ul>
			</div>
	</nav>
	`
};
