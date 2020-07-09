import html from '../html.js';

import ModList from '../ModList.js';
import LicenseBadge from '../LicenseBadge.js';

import state from './state.js';
import { upload } from './actions.js';

export default {
  view: () => html`
    <div class="row">
      <div class="col-md-6">
        <input type="text"
          class="form-control"
          placeholder="Schema name"
          value=${state.name}
          oninput=${e => state.name = e.target.value}
          disabled=${state.progress > 0}
        />
      </div>
      <div class="col-md-6">
        <h3>${state.name}</h3>
      </div>
    </div>
    <div class="row">
      <div class="col-md-6">
        <textarea
          class="form-control"
          style="height: 350px;"
          value=${state.description}
          oninput=${e => state.description = e.target.value}
          disabled=${state.progress > 0}
        />
      </div>
      <div class="col-md-6">
        <pre>${state.description}</pre>
      </div>
    </div>
    <div class="row">
      <div class="col-md-6">
        <input type="text"
          class="form-control"
          placeholder="License"
          value=${state.license}
          oninput=${e => state.license = e.target.value}
          disabled=${state.progress > 0}
        />
      </div>
      <div class="col-md-6">
        <${LicenseBadge} license=${state.license}/>
      </div>
    </div>
    <div class="row">
      <div class="col-md-6">
        <${ModList} schema=${{ mods: state.result.stats}}/>
      </div>
      <div class="col-md-6">
        <ul>
          ${Object.keys(state.result.stats).map(mod => html`<li>${mod}: ${state.result.stats[mod]}</li>`)}
        </ul>
      </div>
    </div>
    <div class="row">
      <div class="col-md-12">
        <button class="btn btn-primary btn-block"
          onclick=${() => upload()}
          disabled=${state.progress > 0}>
          <i class="fa fa-save"></i> Upload
        </button>
      </div>
    </div>
  `
};
