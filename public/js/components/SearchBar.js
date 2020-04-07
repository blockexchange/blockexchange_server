
export default {
  view(vnode) {
    return m("input[type=text]", {
      class: "form-control",
      value: vnode.attrs.keywords,
      placeholder: "Enter search term, for example 'mesecons'",
      oninput: e => vnode.attrs.onChange(e.target.value)
    });
  }
};
