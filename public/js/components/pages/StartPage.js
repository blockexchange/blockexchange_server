import RecentEvents from '../recent/RecentEvents.js';

export default {
	components: {
		"recent-events": RecentEvents
	},
	template: /*html*/`
		<div>
			<div class="text-center">
				<div class="alert alert-danger" role="alert">
					Blockexchange is currently undergoing major refactoring and may be unusable for now
					<br>
					Stay tuned for the next official release announcement
				</div>
				<hr/>
				<h4>Exchange your schemas across worlds with ease</h4>
				<hr/>

				<div>
					<router-link to="/search" class="btn btn-primary">
						<i class="fa fa-search"></i>
						Search
					</router-link>
					<router-link to="/users" class="btn btn-primary">
						<i class="fa fa-users"></i>
						Users
					</router-link>
					<router-link to="/mod" class="btn btn-primary">
						<i class="fa fa-download"></i>
						Mod/Installation
					</router-link>
					<a href="https://github.com/blockexchange" class="btn btn-secondary">
						<i class="fab fa-github"></i>
						Browse the source
					</a>
					<a href="https://discord.gg/ye9aCUJPdB" class="btn btn-secondary">
						<i class="fab fa-discord"></i>
						Join the discord community
					</a>
			</div>
		</div>
		<hr/>

		<recent-events/>
	</div>
	`
};
