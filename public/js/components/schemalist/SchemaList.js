import html from '../html.js';

import LicenseBadge from '../LicenseBadge.js';
import ModList from '../ModList.js';

import prettybytesize from '../../util/prettybytesize.js';

const entry = (schema, hide_user) => html`
	<tr class="${schema.complete ? "" : "table-danger"}">
		${hide_user ? null : html`<td><a href="#!/schema/${schema.user.name}">${schema.user.name}</a></td>`}
		<td>
			<a href="#!/schema/${schema.user.name}/${schema.name}">
				${schema.name}
			</a>
		</td>
		<td>
			${moment(+schema.created).format("YYYY-MM-DD HH:mm")}
			(${moment.duration( moment(+schema.created).diff() ).humanize(true)})
		</td>
		<td>${schema.downloads}</td>
		<td><${LicenseBadge} license=${schema.license}/></td>
		<td>
			<span class="badge badge-secondary">${prettybytesize(schema.total_size)}</span>
		</td>
		<td>${schema.size_x}/${schema.size_y}/${schema.size_z}</td>
		<td>${schema.total_parts}</td>
		<td>${schema.description.substring(0,15)}</td>
		<td><${ModList} schema=${schema}/></td>
	</tr>
`;

export default {
	view: ({ attrs: { list, hide_user }}) => html`
		<table class="table table-striped table-condensed">
			<thead>
				${hide_user ? null : html`<th>User</th>`}
				<th>Name</th>
				<th>Created</th>
				<th>Downloads</th>
				<th>License</th>
				<th>Size [bytes]</th>
				<th>Size [blocks]</th>
				<th>Parts</th>
				<th>Description</th>
				<th>Mods</th>
			</thead>
			<tbody>
				${list.map(e => entry(e, hide_user))}
			</tbody>
		</table>
	`
};
