import LocalLogin from '../login/LocalLogin.js';
import infoStore from '../../store/info.js';

export default {
	components: {
		'local-login': LocalLogin
	},
	data: function(){
		return {
			info: infoStore
		};
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
				<div class="card" v-if="info.enable_signup">
					<div class="card-header">
						Sign up
					</div>
					<div class="card-body">
						<router-link to="/register" class="btn btn-secondary">
							Register
						</router-link>
					</div>
				</div>
				<div class="card">
					<div class="card-header">
						External login
					</div>
					<div class="card-body">
						<a v-if="info.oauth.discord_id"
							v-bind:href="'https://discord.com/api/oauth2/authorize?client_id=' + info.oauth.discord_id + '&redirect_uri=' + encodeURIComponent(info.oauth.base_url + '/api/oauth_callback/discord') + '&response_type=code&scope=identify%20email'"
							class="btn btn-secondary">
							<i class="fab fa-discord"></i>
							Login with Discord
						</a>
						<a v-if="info.oauth.github_id"
							v-bind:href="'https://github.com/login/oauth/authorize?client_id=' + info.oauth.github_id"
							class="btn btn-secondary">
							<i class="fab fa-github"></i>
							Login with Github
						</a>
						<a v-if="info.oauth.mesehub_id"
							v-bind:href="'https://git.minetest.land/login/oauth/authorize?client_id=' + info.oauth.mesehub_id + '&redirect_uri=' + encodeURIComponent(info.oauth.base_url + '/api/oauth_callback/mesehub') + '&response_type=code&state=STATE'"
							class="btn btn-secondary">
							<img src="pics/default_mese_crystal.png"/>
							Login with Mesehub
						</a>
					</div>
				</div>
			</div>
		</div>
	`
};
