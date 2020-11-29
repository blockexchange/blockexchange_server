import NavBar from './components/NavBar.js';

Vue.component('app', {
	components: {
		"nav-bar": NavBar
	},
	template: /*html*/`
		<div>
			<nav-bar/>
			<div class="container-fluid">
				<router-view></router-view>
			</div>
		</div>
	`
});
