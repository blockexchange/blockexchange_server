import { get_all } from '../../api/access_token.js';

export default {
	methods: {
		update: function(){
			get_all()
			.then(l => console.log(l));
		}
	},
	mounted: function(){
		this.update();
	},
	template: /*html*/`
	<div>
		<div class="row">
			<div class="col-md-6">
				<div class="card">
					<div class="card-header">
						Access token
					</div>
					<div class="card-body">
					</div>
				</div>
			</div>
		</div>
	`
};
