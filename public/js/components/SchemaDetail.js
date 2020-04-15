import LicenseBadge from './LicenseBadge.js';
import Preview from './preview/Preview.js';

export default {
	oncreate: function(vnode){
		vnode.state.progress = 0;
	},
  view: function(vnode){
    const schema = vnode.attrs.schema;

    return m("div", [
      m("h3", [
        m("span", { class: "badge badge-primary"}, schema.id),
        " ",
        schema.name,
        " ",
        m("small", { class: "text-muted" }, schema.user.name)
      ]),
      m("div", [
        m("h4", schema.description),
        schema.long_description
      ]),
      m("div", [
        "Size [bytes]: ", m("span", { class: "badge badge-secondary"}, schema.total_size)
      ]),
      m(LicenseBadge, { license: schema.license }),
      m("div", [
        "Created: ",
    		moment(+schema.created).format("YYYY-MM-DD HH:mm"),
    		" (",
    		moment.duration( moment(+schema.created).diff() ).humanize(true),
    		")"
    	]),
			m("span", { class: "badge badge-primary"}, Math.floor(vnode.state.progress) + "%"),
      m(Preview, { schema: schema, progressCallback: f => vnode.state.progress = f * 100 })
    ]);
  }
};
