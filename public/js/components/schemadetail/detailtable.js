import moddetails from './moddetails.js';
import ModList from '../ModList.js';
import LicenseBadge from '../LicenseBadge.js';


export default schema => m("table", { class: "table table-condensed table-striped" }, [
  m("tbody", [
    m("tr", [
      m("td", "Size [bytes]"),
      m("td", m("span", { class: "badge badge-secondary"}, schema.total_size))
    ]),
    m("tr", [
      m("td", "Size [blocks]"),
      m("td", `${schema.size_x} / ${schema.size_y} / ${schema.size_z}`)
    ]),
    m("tr", [
      m("td", "Volume [blocks]"),
      m("td", `${schema.size_x * schema.size_y * schema.size_z}`)
    ]),
    m("tr", [
      m("td", "License"),
      m("td", m(LicenseBadge, { license: schema.license }))
    ]),
    m("tr", [
      m("td", "Created"),
      m("td", [
        moment(+schema.created).format("YYYY-MM-DD HH:mm"),
        " (",
        moment.duration( moment(+schema.created).diff() ).humanize(true),
        ")"
      ])
    ]),
    m("tr", [
      m("td", "Parts"),
      m("td", schema.total_parts)
    ]),
    m("tr", [
      m("td", "Mods"),
      m("td", m(ModList, { schema: schema }))
    ]),
    m("tr", [
      m("td", "Mod block count"),
      m("td", moddetails(schema))
    ])
  ])
]);
