import SearchResult from './SearchResult.js';
import { search } from '../../api/searchschema.js';
import Pager from '../Pager.js';

const store = Vue.reactive({
	result: null,
	total: 0,
	search_params: {
		keywords: ""
	}
});

export default {
	components: {
		"search-result": SearchResult,
		"pager-component": Pager
	},
	data: () => store,
	methods: {
		fetchData: function(limit, offset) {
			const query = {
				order_column: "created",
				order_direction: "desc",
				complete: true
			};

			if (this.search_params.keywords != "") {
				query.keywords = this.search_params.keywords.trim().replaceAll(" ", "|");
			}

			search(query, limit, offset)
			.then(result => {
				this.result = result;
				this.total = result.total;
			});
		},
		search: function() {
			this.$refs.pager.update();
		}
	},
	template: /*html*/`
		<div>
			<form v-on:submit.prevent="search">
				<div class="input-group mb-3">
					<input type="text"
						class="form-control"
						placeholder="Search term (for example 'mesecons')"
						v-model="search_params.keywords"
					/>
					<button class="btn btn-primary" type="button" v-on:click="search">
						<i class="fa fa-search"></i>
						Search
					</button>
				</div>
			</form>
			<pager-component ref="pager" :total="total" v-on:fetchData="fetchData" :limit="20" :route="$route"/>
			<search-result :list="result.list" v-if="result"/>
		</div>
	`
};
