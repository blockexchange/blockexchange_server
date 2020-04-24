
import UserStatus from './UserStatus.js';
import { a, div, ul, li } from './fragments/html.js';
import { fa } from './fragments/fa.js';

import { get_claims } from '../store/token.js';

export default {
  view(){
    return m("nav", { class: "navbar fixed-top navbar-expand-lg navbar-dark bg-dark" }, [
      a({class: "navbar-brand", href: "#!/" }, "Block exchange"),
      div({class:"navbar-collapse collapse"}, [
        ul({class: "navbar-nav mr-auto"}, [
          li({class: "nav-item"},
            a({class:"nav-link",href:"#!/"}, [
              fa("question"),
              "About"
            ])
          ),
          li({class: "nav-item"},
            a({class:"nav-link",href:"#!/login"}, [
              fa("sign-in"),
              "Login"
            ])
          ),
          li({class: "nav-item"},
            a({class:"nav-link",href:"#!/mod"}, [
              fa("download"),
              "Mod/Installation"
            ])
          ),
          li({class: "nav-item"},
            a({class:"nav-link",href:"#!/users"}, [
              fa("users"),
              "Users"
            ])
          ),
          get_claims() ? li({class: "nav-item"},
            a({class:"nav-link",href:"#!/import"}, [
              fa("upload"),
              "Import"
            ])
          ) : null,
          li({class: "nav-item"},
            a({class:"nav-link",href:"#!/search"}, [
              fa("search"),
              "Search"
            ])
          )
        ])
      ]),
      m(UserStatus)
    ]);
  }
};
