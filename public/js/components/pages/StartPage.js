import RecentEvents from '../recent/RecentEvents.js';

export default {
	data: function(){
		return {
			info: null
		};
	},
	components: {
		"recent-events": RecentEvents
	},
	template: /*html*/`
		<div>
			<div class="text-center">
				<h4>Exchange your schemas across worlds with ease</h4>
				<hr/>

				<div>
					<router-link to="/search" class="btn btn-primary">
						<i class="fa fa-search"></i>
						Search
					</router-link>
					&nbsp;
					<router-link to="/users" class="btn btn-primary">
						<i class="fa fa-users"></i>
						Users
					</router-link>
					&nbsp;
					<router-link to="/mod" class="btn btn-primary">
						<i class="fa fa-download"></i>
						Mod/Installation
					</router-link>
					&nbsp;
					<a href="https://github.com/blockexchange" class="btn btn-secondary">
						<i class="fab fa-github"></i>
						Browse the source
					</a>
					&nbsp;
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
