export default {
	template: /*html*/`
	<div>
		<div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8">
				<div class="card">
					<div class="card-header">
						Blockexchange usage
					</div>
					<div class="card-body">
						<h3>Installation</h3>
						<ul>
							<li>
								Download the mod via the "Content" tab in your minetest application or directly from
								<a href="https://github.com/blockexchange/blockexchange">Github</a>
								or the
								<a href="https://content.minetest.net/packages/BuckarooBanzay/blockexchange/">ContentDB</a>
								<img class="rounded" src="pics/help_install_mod.png"/>
							</li>
							<li>
								Enable the secure-http flag for the mod with <b>secure.http_mods = blockexchange</b> in your settings
								<img class="rounded" src="pics/help_install_secure_mods.png"/>
							</li>
						</ul>
						<h3>Download schematics</h3>
						<ul>
							<li>
								Browse and download schemas with <b>/bx_search [name]</b> and <b>/bx_load [user] [name]</b>
								or search for them <router-link to="/search">here</router-link>
							</li>
						</ul>
						<h3>Upload schematics</h3>
						<ul>
							<li>
								To upload schematics you have to <router-link to="/register">register</router-link> and <router-link to="/login">login</router-link> first
								(you can also use external logins for github <i class="fab fa-github"></i>, discord <i class="fab fa-discord"></i> or mesehub <img src="pics/default_mese_crystal.png"/>)
							</li>
							<li>
								After that you can go to your <router-link to="/profile">profile</router-link> page
								and obtain your access-token or create a new one
							</li>
							<li>
								Enter the access-token into your chat-console, for example <b>/bx_login myusername abcdef</b>
								(copy+paste from the browser to the game should work too)
							</li>
							<li>
								Select the area you want to upload with <b>/bx_pos1</b> and <b>/bx_pos2</b>
							</li>
							<li>
								Upload it with <b>/bx_save [schematic name]</b>
							</li>
							<li>
								The schematic should now be visible in the <router-link to="/search">search</router-link> page
							</li>
							<li>
								A reference of the available commands is in the official
								<a href="https://github.com/blockexchange/blockexchange/blob/master/readme.md#chat-commands">mod source</a>
							</li>
						</ul>
					</div>
				</div>
			</div>
			<div class="col-md-2"></div>
	</div>
	`
};
