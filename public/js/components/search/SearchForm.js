import SearchParams from './SearchParams.js';
import SearchResult from './SearchResult.js';
import { search, find_recent } from '../../api/searchschema.js';

export default {
	components: {
		"search-params": SearchParams,
		"search-result": SearchResult
	},
	data: function(){
		return {
			list: []
		};
	},
	mounted: function(){
		find_recent(20).then(l => this.list = l);
	},
	methods: {
		search: function(term){
			search({ keywords: term.trim().replaceAll(" ", "|") })
			.then(l => this.list = l);
		}
	},
	template: /*html*/`
		<div>
			SearchForm
			<search-params v-on:search="search"/>
			<search-result :list="list"/>
		</div>
	`
};
