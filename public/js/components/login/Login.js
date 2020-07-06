import Breadcrumb from '../Breadcrumb.js';
import state from './state.js';

const html = htm.bind(m);

const login_button = () => html`
  <button class="btn btn-primary btn-block"
    disabled=${!state.username || !state.password}
    onclick=${state.login}>
      Login
      <span class="badge badge-danger">${state.message}</span>
  </button>
`;

const temp_login_button = () => html`
  <button class="btn btn-secondary btn-block"
    onclick=${state.temp_login}>
    Login with a temporary account
  </button>
`;

const logout_button = () => html`
  <button class="btn btn-primary btn-block"
    onclick=${state.logout}>
    Logout
  </button>
`;

const breadcrumb_links = [
  { name: "Home", link: "#!/" },
  { name: "Login", active: true }
];

export default {
  view: () => html`
    <${Breadcrumb} links=${breadcrumb_links}/>
    <div class="row">
    <div class="col-md-4"></div>
    <div class="col-md-4">
      <form>
        <input type="text"
          class="form-control"
          placeholder="Username"
          value=${state.username}
          oninput=${e => state.username = e.target.value}/>
        <input type="password"
          class="form-control"
          placeholder="Password"
          value=${state.password}
          oninput=${e => state.password = e.target.value}/>
        ${state.isLoggedIn() ? logout_button() : login_button()}
        ${state.isLoggedIn() ? null : temp_login_button()}
      </form>
    </div>
    <div class="col-md-4"></div>
    </div>
  `
};
