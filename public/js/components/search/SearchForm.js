import SearchResult from './SearchResult.js';
import { search } from '../../api/searchschema.js';
import { get_all as get_tags } from '../../api/tag.js';
import Pager from '../Pager.js';

const store = Vue.reactive({
	result: null,
	total: 0,
	tags: [],
	selected_tag: null,
	keywords: ""
});

export default {
	components: {
		"search-result": SearchResult,
		"pager-component": Pager
	},
	created: function() {
		if (this.tags.length == 0){
			// fetch tags
			get_tags().then(t => this.tags = t);
		}
	},
	data: () => store,
	methods: {
		fetchData: function(limit, offset) {
			const query = {
				order_column: "created",
				order_direction: "desc",
				complete: true
			};

			if (this.selected_tag) {
				query.tag_id = this.selected_tag.id;
			}
			
			if (this.keywords != "") {
				query.keywords = this.keywords.trim().replaceAll(" ", "|");
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
				<div class="row">
					<div class="col-md-8">
						<input type="text"
							class="form-control"
							placeholder="Search term (for example 'mesecons')"
							v-model="keywords"
						/>
					</div>
					<div class="col-md-3">
						<div class="input-group">
							<span class="input-group-text">Tag</span>
							<select class="form-select" v-model="selected_tag">
								<option></option>
								<option v-for="tag in tags" :value="tag">
									{{ tag.name }}
								</option>
							</select>
						</div>
					</div>
					<div class="col-md-1">
						<button class="btn btn-primary w-100" type="button" v-on:click="search">
							<i class="fa fa-search"></i>
							Search
						</button>
					</div>
				</div>
			</form>
			<pager-component ref="pager" :total="total" v-on:fetchData="fetchData" :limit="20" :route="$route"/>
			<search-result :list="result.list" v-if="result"/>
		</div>
	`
};
