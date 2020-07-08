import html from './html.js';

export default {
  view: ({ attrs: { schema }}) => html`
    <div class="card">
      <div class="card-body">
        <p>Usage:</p>
        <p>/bx_pos1</p>
        <p>/bx_load ${schema.user.name} ${schema.name}</p>
      </div>
    </div>
  `
};
