import html from './html.js';

export default {
	view: function(){
		const route = m.route.get();
		if (!route){
			return;
		}

		const parts = route.split("/");
		const items = [];

		items.push(html`
			<li class="breadcrumb-item">
				<a href="#!/">Home</a>
			</li>
		`);

		if (route == "/search") {
			items.push(html`<li class="breadcrumb-item">Search</li>`);
		}

		if (route == "/login") {
			items.push(html`<li class="breadcrumb-item">Login</li>`);
		}

		if (route == "/users") {
			items.push(html`<li class="breadcrumb-item">Users</li>`);
		}

		if (route == "/register") {
			items.push(html`<li class="breadcrumb-item">Register</li>`);
		}

		if (route == "/mod") {
			items.push(html`<li class="breadcrumb-item">Mod</li>`);
		}

		if (route == "/import") {
			items.push(html`<li class="breadcrumb-item">Import</li>`);
		}

		if (parts[1] == "schema"){
			items.push(html`<li class="breadcrumb-item">User-Schemas</li>`);

			let link = "#!/schema";

			if (parts.length >= 3){
				link += "/" + parts[2];
				items.push(html`<li class="breadcrumb-item">
					<a href="${link}">${parts[2]}</a>
				</li>`);
			}

			if (parts.length >= 4){
				let sane_schemaname = parts[3].replaceAll("%20", " ");
				link += "/" + sane_schemaname;
				items.push(html`<li class="breadcrumb-item">
					<a href="${link}">${sane_schemaname}</a>
				</li>`);
			}

			if (parts.length >= 5){
				if (parts[4] == "edit"){
					items.push(html`<li class="breadcrumb-item">Edit</li>`);
				} else if (parts[4] == "preview"){
					items.push(html`<li class="breadcrumb-item">Preview</li>`);
				}
			}
		}

		return items;

	}
};
