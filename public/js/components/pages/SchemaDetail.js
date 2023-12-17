import Breadcrumb, { START, SCHEMA_DETAIL } from "../Breadcrumb.js";
import LoadingBlock from "../LoadingBlock.js";

import format_time from "../../util/format_time.js";
import format_size from "../../util/format_size.js";
import { get_schema_by_name } from "../../api/schema.js";

export default {
    props: ["username", "name"],
	components: {
        "bread-crumb": Breadcrumb,
		"loading-block": LoadingBlock
	},
	data: function() {
		return {
			breadcrumb: [START, SCHEMA_DETAIL(this.username, this.name)],
			BaseURL
		};
	},
	methods: {
		format_time,
		format_size,
		fetch_data: function() {
			return {
				schema: get_schema_by_name(this.username, this.name)
			};
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<loading-block :fetch_data="fetch_data" v-slot="{ data }">
			{{data.schema.description}}
			<div class="row">
				<div class="col-md-6">
					<h3>
						{{data.schema.name}}
						<small class="text-muted">by {{username}}</small>
						&nbsp;
						<button class="btn btn-outline-primary">
							<i class="fa fa-star" v-bind:style="{color: 'yellow'}"></i>
							<span class="badge bg-secondary rouded-pill">{{data.schema.stars}}</span>
							Star
						</button>
					</h3>					
				</div>
				<div class="col-md-6">
					<div class="btn-group float-end">
						<a class="btn btn-sm btn-secondary">
							<i class="fa fa-edit"></i> Edit
						</a>
						<a class="btn btn-sm btn-secondary">
							<i class="fa fa-image"></i> Update screenshot
						</a>
						<a class="btn btn-sm btn-danger">
							<i class="fa fa-trash"></i> Delete
						</a>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-md-4">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Details</h5>
							<ul>
								<li><b>Created:</b> {{format_time(data.schema.created / 1000)}}</li>
								<li><b>Modified:</b> {{format_time(data.schema.mtime / 1000)}}</li>
								<li><b>Size:</b> {{format_size(data.schema.total_size)}}</li>
								<li><b>Dimensions:</b> {{data.schema.size_x}} / {{data.schema.size_y}} / {{data.schema.size_z}} nodes</li>
								<li><b>Parts:</b> {{data.schema.total_parts}}</li>
								<li><b>Downloads:</b> {{data.schema.downloads}}</li>
								<li>
									<b>License:</b>
									<img v-if="data.schema.license == 'CC0'" :src="BaseURL + '/pics/license_cc0.png'">
									<img v-else-if="data.schema.license == 'CC-BY-SA'" :src="BaseURL + '/pics/license_cc-by-sa.png'">
									<span v-else class="badge bg-secondary">{{data.schema.license}}</span>
								</li>
							</ul>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Description</h5>
						</div>
					</div>
				</div>
			</div>
		</loading-block>
	`
};
