import LicenseBadge from './LicenseBadge.js';

export default {
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
    ]);
  }
};
