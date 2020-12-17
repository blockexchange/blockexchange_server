import UpdateProfile from './UpdateProfile.js';
import store from '../../store/login.js';

export default {
	components: {
		"update-profile": UpdateProfile
	},
	data: function(){
		return {
			store: store
		};
	},
	template: /*html*/`
		<div class="row">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						Update profile
					</div>
					<div class="card-body">
						<update-profile/>
					</div>
				</div>
			</div>
			<div class="col-md-6">
				<div class="card">
				  <div class="card-header">
				    Profile info
				  </div>
				  <div class="card-body">
						<ul>
							<li>
								ID: <b>{{ store.claims.user_id }}</b>
							</li>
							<li>
								Name: <b>{{ store.claims.username }}</b>
							</li>
							<li>
								Mail: <b>{{ store.claims.mail }}</b>
							</li>
							<li>
								Role: <span class="badge badge-success">{{ store.claims.role }}</span>
							</li>
							<li>
								Type: <span class="badge badge-success">{{ store.claims.type }}</span>
							</li>
						</ul>
				  </div>
				</div>
			</div>
		</div>
	`
};
