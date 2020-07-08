
import UserStatus from './UserStatus.js';
import html from './html.js';

import { get_claims } from '../store/token.js';

const logged_in_nav = () => html`
  <li class="nav-item">
    <a class="nav-link" href="#!/import">
      <i class="fa fa-upload"></i> Import
    </a>
  </li>
`;

export default {
  view: () => html`
    <nav class="navbar fixed-top navbar-expand-lg navbar-dark bg-dark">
      <a class="navbar-brand" href="#!/">Block exchange</a>
      <div class="navbar-collapse collapse">
        <ul class="navbar-nav mr-auto">
          <li class="nav-item">
            <a class="nav-link" href="#!/">
              <i class="fa fa-question"></i> About
            </a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#!/login">
              <i class="fa fa-sign-in"></i> Login
            </a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#!/mod">
              <i class="fa fa-download"></i> Mod/Installation
            </a>
          </li>
          <li class="nav-item">
            <a class="nav-link" href="#!/users">
              <i class="fa fa-users"></i> Users
            </a>
          </li>
          ${get_claims() ? logged_in_nav() : null}
          <li class="nav-item">
            <a class="nav-link" href="#!/search">
              <i class="fa fa-search"></i> Search
            </a>
          </li>
        </ul>
      </div>
      <${UserStatus}/>
    </nav>
  `
};
