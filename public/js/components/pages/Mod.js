import Breadcrumb, { START, MOD } from "../Breadcrumb.js";
import ClipboardCopy from "../ClipboardCopy.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"clipboard-copy": ClipboardCopy
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
				<h5>Installation instructions (youtube)</h5>
				<hr>
				<iframe width="420" height="315" src="https://www.youtube.com/embed/x0geU_EyO-0"></iframe>
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
				This example uses the ingame upload mechanism, you can use the web-based file-<router-link to="/import">import</router-link> as an alternative
				<ul>
					<li>Obtain an access-token from your <router-link to="/profile">profile</router-link> page</li>
					<li>You must have either have the <b>blockexchange</b> or the <b>blockexchange_protected_upload</b> privilege to be able to upload</li>
					<li>Enter the access-token command (<b>/bx_login</b>) in a server with blockexchange installed and enabled</li>
					<li>Mark the area you want to upload with <clipboard-copy :text="'/bx_pos1'"></clipboard-copy> and <clipboard-copy :text="'/bx_pos2'"></clipboard-copy></li>
					<li>Upload it with <clipboard-copy text="/bx_save [schematic name]"></clipboard-copy></li>
					<li>The schematic should now be visible in the <router-link to="/search">search</router-link> page</li>
					<li>A reference of the available commands is in the official <a href="https://github.com/blockexchange/blockexchange/blob/master/readme.md#chat-commands">mod source</a></li>
				</ul>
			</div>
			<div class="col-md-2"></div>
		</div>
	`
};
