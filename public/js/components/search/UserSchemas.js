import SearchResult from './SearchResult.js';
import { search } from '../../api/searchschema.js';
import Pager from '../Pager.js';

export default {
    props: ["user_name"],
	components: {
		"search-result": SearchResult,
		"pager-component": Pager
	},
	data: function() {
        return {
            total: 0,
            result: null
        };
    },
	methods: {
		fetchData: function(limit, offset) {
			const query = {
				order_column: "created",
				order_direction: "desc",
				complete: true,
                user_name: this.user_name
			};

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
			<pager-component ref="pager" :total="total" v-on:fetchData="fetchData" :limit="20" :route="$route"/>
			<search-result :list="result.list" v-if="result"/>
		</div>
	`
};
