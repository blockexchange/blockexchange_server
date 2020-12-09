export default {
	data: function(){
		return {
			username: "",
			password: "",
			password2: "",
			mail: ""
		};
	},
	template: /*html*/`
		<form v-on:submit.prevent>
			<input type="text"
				class="form-control"
				placeholder="Username"
				v-model="username"
			/>
			<input type="text"
				class="form-control"
				placeholder="E-Mail"
				v-model="mail"
			/>
			<input type="password"
				class="form-control"
				placeholder="Password"
				v-model="password"
			/>
			<input type="password"
				class="form-control"
				placeholder="Password (verify)"
				v-model="password2"
			/>
			<button class="btn btn-secondary btn-block">
				Register
			</button>
		</form>
	`
};
