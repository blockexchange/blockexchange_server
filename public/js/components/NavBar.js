import NavDropdown from "./NavDropdown.js";

export default {
	components: {
		"nav-dropdown": NavDropdown
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand">Blockexchange</router-link>
				<ul class="navbar-nav me-auto mb-2 mb-lg-0">
					<li class="nav-item">
						<router-link to="/" class="nav-link">
							<i class="fa fa-home"></i> Start
						</router-link>
					</li>
				</ul>
			</div>
		</nav>
	`
};
