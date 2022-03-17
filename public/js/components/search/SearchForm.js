import SearchParams from './SearchParams.js';
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
		"search-params": SearchParams,
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
		}
	},
	template: /*html*/`
		<div>
			<search-params v-on:search="this.$refs.pager.update()" :search_params="search_params"/>
			<pager-component ref="pager" :total="total" v-on:fetchData="fetchData" :limit="20" :route="$route"/>
			<search-result :list="result.list" v-if="result"/>
		</div>
	`
};
