import html from '../html.js';

export default schema => html`
	<ul>
		<div>
			${Object.keys(schema.mods).map(mod => html`<li>${mod}: ${schema.mods[mod]}</li>`)}
		</div>
	</ul>
`;
