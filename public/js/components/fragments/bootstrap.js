
export const button = (type, link, content) => m("a", {
  class: `btn btn-${type}`,
  href: link
}, content);

export const button_secondary = (link, content) => button("secondary", link, content);
