import SearchParams from './SearchParams.js';
import SearchResult from './SearchResult.js';
import { search } from '../../api/searchschema.js';

const store = Vue.reactive({
	list: [],
	term: ""
});

export default {
	components: {
		"search-params": SearchParams,
		"search-result": SearchResult
	},
	data: () => store,
	mounted: function(){
		this.search(this.term);
	},
	methods: {
		search: function(term){
			this.term = term.trim().replaceAll(" ", "|");
			if (this.term === "") {
				// initialize list with recent additions
				const q = {
					order_column: "created",
					order_direction: "desc",
					complete: true
				};
				search(q, 20, 0).then(l => this.list = l);
			} else {
				search({ keywords: this.term, complete: true }, 20, 0)
				.then(l => this.list = l);
			}

		}
	},
	template: /*html*/`
		<div>
			<search-params v-on:search="search" :term="term"/>
			<search-result :list="list"/>
		</div>
	`
};
