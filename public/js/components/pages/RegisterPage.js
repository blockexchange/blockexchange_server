import RegisterForm from '../login/RegisterForm.js';

export default {
	components: {
		"register-form": RegisterForm
	},
	template: /*html*/`
		<div>
			<div class="row">
				<div class="col-md-2"></div>
				<div class="col-md-8">
					<div class="card">
						<div class="card-header">
							Register
						</div>
						<div class="card-body">
							<register-form/>
						</div>
					</div>
				</div>
				<div class="col-md-2"></div>
		</div>
	`
};
