import Breadcrumb, { START, MOD } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, MOD]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8">
				<h5>Installation instructions</h5>
				<hr>
				<p>Download the mod via "Content" tab in your minetest application</p>
				<img src="pics/help_install_mod.png" class="rounded"/>
				<hr>
				<p>
					Enable the secure-http setting for the mod with
					<b class="text-muted">secure.http_mods = blockexchange</b>
					in your settings
				</p>
				<img src="pics/help_install_secure_mods.png" class="rounded"/>
				<br>
				<br>
				<div class="alert alert-secondary">
					<i class="fa fa-circle-info"></i>
					<b>Note:</b> The "online" download only works if the setting is set properly
				</div>
				<hr>
				<h5>Download</h5>
				<p>
					<router-link to="/search">Browse</router-link> schematics and follow the download instructions on the detail page
				</p>
				<hr>
				<h5>Upload</h5>
				<div class="alert alert-secondary">
					<i class="fa fa-circle-info"></i>
					<b>Note:</b> You need to be logged in in order to upload schematics
				</div>
				<p>TODO</p>
			</div>
			<div class="col-md-2"></div>
		</div>
	`
};
