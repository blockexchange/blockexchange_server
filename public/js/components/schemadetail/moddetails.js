
export default schema => m("ul",
	m("div", Object.keys(schema.mods).map(
		mod_name => m("li", `${mod_name}: ${schema.mods[mod_name]}`)
	)));
