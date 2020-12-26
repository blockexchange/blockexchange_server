export default {
	template: /*html*/`
	<div>
		<div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8">
				<div class="card">
					<div class="card-header">
						Mod installation
					</div>
					<div class="card-body">
						<ul>
							<li>
								Download the mod via the "Content" tab in your minetest application or directly from
								<a href="https://github.com/blockexchange/blockexchange">Github</a>
								or the
								<a href="https://content.minetest.net/packages/BuckarooBanzay/blockexchange/">ContentDB</a>
							</li>
							<li>
								Enable the secure-http flag for the mod with <b>secure.http_mods = blockexchange</b> in your settings
							</li>
							<li>
								Browse and download schemas with <b>/bx_search [name]</b> and <b>/bx_load [user] [name]</b>
								or search for them <router-link to="/search">here</router-link>
							</li>
						</ul>
					</div>
				</div>
			</div>
			<div class="col-md-2"></div>
	</div>
	`
};
