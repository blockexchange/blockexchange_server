import Breadcrumb from './Breadcrumb.js';
import html from './html.js';

const links = [{
  name: "Home",
  link: "#!/"
},{
  name: "Mod",
  active: true
}];

export default {
  view: () => html`
    <${Breadcrumb} links=${links}/>
    <h4>Mod/Installation</h4>
    <p>Check out the readme on the mod repository:
      <a href="https://github.com/blockexchange/blockexchange#basic-usage">
        https://github.com/blockexchange/blockexchange
      </a>
    </p>
  `
};
