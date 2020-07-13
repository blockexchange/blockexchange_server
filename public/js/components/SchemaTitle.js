import html from './html.js';

export default {
  view: ({ attrs: { schema }}) => html`
    <h3 style="display: inline;">
      <span class="badge badge-primary">${schema.id}</span>
      ${schema.name}
      <small class="text-muted"> by ${schema.schemagroup.name}</small>
      ${schema.complete ? "" : html`<span class="badge badge-warning>Incomplete!</span>`}
    </h3>
  `
};
