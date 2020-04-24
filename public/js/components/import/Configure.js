
import { row, col6, col12 } from '../fragments/bootstrap.js';
import { hr, textarea, pre } from '../fragments/html.js';
import { fa } from '../fragments/fa.js';

import ModList from '../ModList.js';
import LicenseBadge from '../LicenseBadge.js';

import state from './state.js';
import { upload } from './actions.js';

export default {
  view: function(){
    return [
      row([
        col6([
          m("input", {
            class: "form-control",
            placeholder: "Schema name",
            value: state.name,
            oninput: e => state.name = e.target.value,
            disabled: state.progress > 0
          })
        ]),
        col6([
          m("h3", state.name)
        ])
      ]),
      row([
        col6([
          textarea({
            value: state.description,
            class: "form-control",
            style: "height: 350px",
            oninput: e => state.description = e.target.value,
            disabled: state.progress > 0
          })
        ]),
        col6([
          pre(state.description)
        ])
      ]),
      row([
        col6([
          m("input", {
            class: "form-control",
            placeholder: "License",
            value: state.license,
            oninput: e => state.license = e.target.value,
            disabled: state.progress > 0
          })
        ]),
        col6([
          m(LicenseBadge, { license: state.license })
        ])
      ]),
      row([
        col6(
          m(ModList, { schema: { mods: state.result.stats }})
        ),
        col6(
          m("ul",
          	m("div", Object.keys(state.result.stats).map(
          		mod_name => m("li", `${mod_name}: ${state.result.stats[mod_name]}`)
          	))
          )
        )
      ]),
      hr(),
      row(col12([
        m("button", {
          class: "btn btn-primary btn-block",
          onclick: () => upload(),
          disabled: state.progress > 0
        },[
          fa("save"),
          " Upload"
        ])
      ]))
    ];
  }
};
