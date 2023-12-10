import Breadcrumb, { START, REGISTER } from "../Breadcrumb.js";
import { create_captcha } from "../../api/captcha.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	mounted: function() {
		create_captcha().then(c => {
			this.captcha_src = `api/captcha/${c}.png`;
		});
	},
	data: function() {
		return {
			username: "",
			password: "",
			password2: "",
			captcha_answer: "",
			captcha_src: null,
			breadcrumb: [START, REGISTER]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
				<h5>Register a new account</h5>
				<input type="text" v-model="username" class="form-control w-100" placeholder="Username"/>
				<input type="password" v-model="password" class="form-control w-100" placeholder="Password"/>
				<input type="password" v-model="password2" class="form-control w-100" placeholder="Retype password"/>
				<img :src="captcha_src" v-if="captcha_src">
				<div class="alert alert-secondary" v-if="!captcha_src">
					<i class="fa fa-spinner fa-spin"></i>
					Loading captcha
				</div>
				<input type="text" v-model="captcha_answer" class="form-control w-100" placeholder="Captcha answer"/>
				<button class="btn btn-primary w-100">
					<i class="fa fa-user-plus"></i>
					Register
				</button>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
