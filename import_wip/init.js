const luaparse = require("./luaparse");
const fs = require("fs");

const schema = fs.readFileSync('./plain_chest.we');

const cleanlua = schema.toString().substring(2);
const ast = luaparse.parse(cleanlua);

console.log(JSON.stringify(ast.body[0].arguments[0].fields));
