import { is_logged_in, get_claims, has_permission } from "../service/login.js";

export default {
	computed: {
		is_logged_in,
		claims: get_claims
	},
	methods: {
		has_permission
	},
	template: /*html*/`
		<nav class="navbar navbar-expand-lg navbar-dark bg-dark">
			<div class="container-fluid">
				<router-link to="/" class="navbar-brand d-none d-lg-inline">Block exchange</router-link>
				<ul class="navbar-nav me-auto">
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
					<li class="nav-item" v-if="is_logged_in">
						<router-link :to="'/schema/' + claims.username" class="nav-link">
							<i class="fa fa-home"></i> My schematics
						</router-link>
					</li>
					<li class="nav-item" v-if="is_logged_in">
						<router-link to="/import" class="nav-link">
							<i class="fa fa-upload"></i> Schematic import
						</router-link>
					</li>
					<li class="nav-item" v-if="is_logged_in">
						<router-link to="/tags" class="nav-link">
							<i class="fa fa-tags"></i> Tags
						</router-link>
					</li>
				</ul>
				<form class="d-flex" v-if="is_logged_in && has_permission('ADMIN')">
					<div class="btn btn-secondary">
						<router-link to="/profile">
							<i class="fas fa-user"></i>
							<span>
								Logged in as <b>{{claims.username}}</b>
							</span>
						</router-link>
					</div>
				</form>
			</div>
		</nav>
	`
};
