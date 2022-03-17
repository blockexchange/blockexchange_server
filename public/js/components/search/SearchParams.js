
export default {
	props: ["search_params"],
	methods: {
		search: function(){
			this.$emit("search");
			return false;
		}
	},
	template: /*html*/`
		<form v-on:submit="search">
			<div class="input-group mb-3">
				<input type="text"
					class="form-control"
					placeholder="Search term (for example 'mesecons')"
					v-model="search_params.keywords"
				/>
				<button class="btn btn-primary" type="button" v-on:click="search">
					Search
				</button>
			</div>
		</form>
	`
};
