import users from '../service/users.js';
import Breadcrumb from './Breadcrumb.js';

import html from './html.js';

const links = [{
  name: "Home",
  link: "#!/"
},{
  name: "Users",
  active: true
}];

const UserEntry = user => html`
<li>
  <a href="#!/schema/${user.name}">
    ${user.name}
    ${" "}
    (created ${moment.duration( moment(+user.created).diff() ).humanize(true)})
  </a>
</li>
`;

var list;

export default {
  oninit: function(){
    if (!list){
      users().then(u => list = u);
      list = [];
    }
  },
  view: () => html`
    <${Breadcrumb} links=${links}/>
    <ul>
      ${list.map(UserEntry)}
    </ul>
  `
};
