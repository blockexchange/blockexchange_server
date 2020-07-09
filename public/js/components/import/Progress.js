import state from './state.js';
import html from '../html.js';

export default {
  view: () => html`
    <div class="progress">
      <div class="progress-bar" style="width: ${state.progress}%">
        ${Math.floor(state.progress * 10) / 10}%
      </div>
    </div>
  `
};
