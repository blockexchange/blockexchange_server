import html from './html.js';

const badge = (cl, txt) => html`
	<span class="badge badge-${cl}">${txt}</span>
`;

const get_badge_class = mod_name => {
	switch(mod_name) {
		case "ignore":
			return "danger";
		case "default":
			return "success";
		default:
			return "secondary";
	}
};

export default {
	view: ({ attrs: { schema }}) => html`
		<div>
			${Object.keys(schema.mods).map(mod_name => badge(get_badge_class(mod_name), mod_name))}
		</div>
	`
};
