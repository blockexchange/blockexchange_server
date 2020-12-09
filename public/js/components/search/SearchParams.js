import debounce from '../../util/debounce.js';

export default {
	methods: {
		search: debounce(function(term){
			this.$emit("search", term);
		}, 250)
	},
	template: /*html*/`
		<form v-on:submit.prevent>
			<input type="text"
				class="form-control"
				placeholder="Search term (for example 'mesecons')"
				v-on:input="search($event.target.value)"
			/>
		</form>
	`
};
