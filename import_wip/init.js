const luaparse = require("./luaparse");
const fs = require("fs");

const schema = fs.readFileSync('./blockexchange.we');

const cleanlua = schema.toString().substring(2);
const ast = luaparse.parse(cleanlua);

//console.log(ast.body[0].arguments[0].fields[0].value.fields);
//console.log(ast.body[0].arguments[0])


function unbox(field){
	if (field.type == "NumericLiteral")
		return +field.raw;
	if (field.type == "StringLiteral")
		return field.raw.substring(1, field.raw.length-1);
}

ast.body[0].arguments[0].fields.forEach(function(o){
	//console.log(o.value.fields);

	const entry = {};

	o.value.fields
	.filter(f => f.type == "TableKey")
	.forEach(function(f){
		const key = unbox(f.key);
		const value = unbox(f.value);
		entry[key] = value;
	});

	console.log(entry);
});
