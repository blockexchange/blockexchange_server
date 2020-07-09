import html from '../html.js';

import moddetails from './moddetails.js';
import ModList from '../ModList.js';
import LicenseBadge from '../LicenseBadge.js';
import prettybytesize from '../../util/prettybytesize.js';

export default schema => html`
  <table class="table table-condensed table-striped">
    <tr>
      <td>Size [bytes]</td>
      <td>
        <span class="badge badge-secondary">${prettybytesize(schema.total_size)}</span>
      </td>
    </tr>
    <tr>
      <td>Size [blocks]</td>
      <td>${schema.size_x} / ${schema.size_y} / ${schema.size_z}</td>
    </tr>
    <tr>
      <td>Volume [blocks]</td>
      <td>${schema.size_x * schema.size_y * schema.size_z}</td>
    </tr>
    <tr>
      <td>License</td>
      <td>
        <${LicenseBadge} license=${schema.license}/>
      </td>
    </tr>
    <tr>
      <td>Created</td>
      <td>
        ${moment(+schema.created).format("YYYY-MM-DD HH:mm")}
        (${moment.duration( moment(+schema.created).diff() ).humanize(true)})
      </td>
    </tr>
    <tr>
      <td>Parts</td>
      <td>${schema.total_parts}</td>
    </tr>
    <tr>
      <td>Mods</td>
      <td>
        <${ModList} schema=${schema}/>
      </td>
    </tr>
    <tr>
      <td>Mod block count</td>
      <td>${moddetails(schema)}</td>
    </tr>
  </table>
`;
