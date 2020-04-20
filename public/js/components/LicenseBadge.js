
export default {
  view: function(vnode){
    switch (vnode.attrs.license) {
      case "CC0":
  			return m("img", {src:"pics/license_cc0.png"});
      case "CC-BY":
  			return m("img", {src:"pics/license_cc-by.png"});
      case "CC-BY-SA":
  			return m("img", {src:"pics/license_cc-by-sa.png"});
      case "CC-BY-NC":
  			return m("img", {src:"pics/license_cc-by-nc.png"});
  		default:
  			return m("span", {class: "badge badge-primary"}, vnode.attrs.license);
  	}
  }
};
