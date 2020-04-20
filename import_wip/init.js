const luaparse = require("./luaparse");
const fs = require("fs");

const schema = fs.readFileSync('./blockexchange.we');

const cleanlua = schema.toString().substring(2);
const ast = luaparse.parse(cleanlua);

console.log(ast.body[0].arguments[0].fields[0].value.fields);
