

function unbox(field){
  if (!field){
    return;
  }

  switch (field.type){
    case "NumericLiteral":
  		return +field.raw;

  	case "StringLiteral":
  		return field.raw.substring(1, field.raw.length-1);

    case "TableConstructorExpression":
      if (field.fields.length > 0){
        if (field.fields[0].type == "TableKey"){
          //Map
          const result = {};
          field.fields.forEach(function(child){
            const key = unbox(child.key);
            const value = unbox(child.value);
            result[key] = value;
          });

          return result;

        } else if (field.fields[0].type == "TableValue"){
          //Array
          const result = [];
          field.fields.forEach(function(child){
            const value = unbox(child.value);
            result.push(value);
          });

          return result;
        } else {
          throw new Error("TableConstructorExpression with unknown child!");
        }
      }
      // type not known
      return {};

    case "TableValue":
      return unbox(field.value);

    default:
      console.error(field);
      throw new Error("Unknown type: " + field.type);
  }
}

export default unbox;
