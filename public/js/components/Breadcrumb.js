
import html from './html.js';

export default {
	view: function(vnode){
		const links = vnode.attrs.links.map(l => html`
			<li class="breadcrumb-item ${l.active ? "active" : ""}">
				${l.link ? html`<a href=${l.link}>${l.name}</a>` : l.name}
			</li>
		`);

		return html`
			<nav>
				<ol class="breadcrumb">
					${links}
				</ol>
			</nav>
		`;
	}
};
