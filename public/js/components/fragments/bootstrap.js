
export const button = (type, link, content) => m("a", {
  class: `btn btn-${type}`,
  href: link
}, content);

export const button_secondary = (link, content) => button("secondary", link, content);

export const row = content => m("div", { class: "row" }, content);

export const col1 = content => m("div", { class: "col-md-1" }, content);
export const col2 = content => m("div", { class: "col-md-2" }, content);
export const col3 = content => m("div", { class: "col-md-3" }, content);
export const col4 = content => m("div", { class: "col-md-4" }, content);
export const col5 = content => m("div", { class: "col-md-5" }, content);
export const col6 = content => m("div", { class: "col-md-6" }, content);
export const col10 = content => m("div", { class: "col-md-10" }, content);
export const col12 = content => m("div", { class: "col-md-12" }, content);
