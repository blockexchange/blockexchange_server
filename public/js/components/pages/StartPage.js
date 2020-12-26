export default {
	template: /*html*/`
		<div class="text-center">
			<img class="img-fluid rounded" src="pics/blockexchange.png"/>
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
				Source
				</a>
		</div>
	</div>
	`
};
