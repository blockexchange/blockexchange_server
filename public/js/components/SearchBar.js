import html from './html.js';

export default {
  view: ({ attrs: { keywords, onChange }}) => html`
    <input type="text"
      class="form-control"
      value=${keywords}
      placeholder="Enter search term, for example 'mesecons'"
      oninput=${e => onChange(e.target.value)}
    />
  `
};
