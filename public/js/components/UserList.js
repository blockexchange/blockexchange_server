import users from '../service/users.js';
import Breadcrumb from './Breadcrumb.js';

var list;

export default {
  oninit: function(){
    if (!list){
      users().then(u => list = u);
      list = [];
    }
  },
  view: function(){
    return [
			m(Breadcrumb, {
				links: [{
					name: "Home",
					link: "#!/"
				},{
					name: "Users",
					active: true
				}]
			}),
			m("ul", list.map(user => {
	      return m("li", [
	        m("a", { href: `#!/schema/${user.name}` }, user.name),
	        " (Created ",
	        moment.duration( moment(+user.created).diff() ).humanize(true),
	        ")"
	      ]);
	    }))
		];
  }
};
