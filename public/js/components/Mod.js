import Breadcrumb from './Breadcrumb.js';


export default {
  view: function() {
    return [
			m(Breadcrumb, {
				links: [{
					name: "Home",
					link: "#!/"
				},{
					name: "Mod",
					active: true
				}]
			}),
			m.trust(`
	      <h4>Mod/Installation</h4>
	      <p>Check out the readme on the mod repository:
	        <a href="https://github.com/blockexchange/blockexchange#basic-usage">
	          https://github.com/blockexchange/blockexchange
	        </a>
	      </p>

	    `)
		];
  }
};
