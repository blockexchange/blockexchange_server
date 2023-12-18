import Breadcrumb, { START, SCHEMA_IMPORT } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SCHEMA_IMPORT]
		};
	},
	methods: {
		on_upload: function() {
			console.log(this.$refs.upload.files);
			this.$refs.upload.files[0].arrayBuffer()
			.then(buf => console.log(buf));
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4">
				<h5>
					Upload Worldedit or Blockexchange schematics
				</h5>
				<div class="input-group">
					<input ref="upload" type="file" class="form-control" v-on:change="on_upload"/>
					<button class="btn btn-primary">
						<i class="fa fa-upload"></i>
						Import
					</button>
				</div>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
