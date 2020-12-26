import { create } from '../../api/token.js';


export default {
	data: function(){
		return {
			permission: "UPLOAD_ONLY",
			expiration: 604800,
			token: null
		};
	},
	methods: {
		updateToken: function(){
			create({
				upload_only: this.permission === "UPLOAD_ONLY",
				expiresIn: +this.expiration
			}).then(token => {
				this.token = token;
			});
		}
	},
	mounted: function(){
		this.updateToken();
	},
	watch: {
		permission: function(){
			this.updateToken();
		},
		expiration: function(){
			this.updateToken();
		}
	},
	template: /*html*/`
		<div>
			<div class="row">
				<div class="col-md-6">
					<div class="card">
						<div class="card-header">
							Create token
						</div>
						<div class="card-body">
							Permissions:
							<select class="form-control" v-model="permission">
								<option value="UPLOAD_ONLY">Upload only</option>
							</select>
							Expiration:
							<select class="form-control" v-model="expiration">
								<option value="86400">1 day</option>
								<option value="604800">1 week</option>
								<option value="2419200">1 month</option>
								<option value="31536000">1 year</option>
							</select>
							<br>
							<pre v-if="token">/bx_token {{ token }}</pre>
						</div>
					</div>
				</div>
		</div>
	`
};
