
export default {
  view: function(vnode){
    switch (vnode.attrs.license) {
  		case "CC0":
  			return m("img", {src:"pics/license_cc0.png"});
  		default:
  			return m("span", {class: "badge badge-primary"}, vnode.attrs.license);
  	}
  }
};
