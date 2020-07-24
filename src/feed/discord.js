const events = require("../eevnts");
const axios = require('axios');

module.exports = function(){
  // listen to events and dispatch webhook

  const schema_feed_url = process.env.DISCORD_SCHEMA_FEED_URL;

  if (schema_feed_url){
    events.on("new-schema", function(schema){
      //https://birdie0.github.io/discord-webhooks-guide/examples/spotify.html
      const data = {
        content: `Schema created: **${schema.name}**
        License: **${schema.license}**
        Size: ${schema.size_x}/${schema.size_y}/${schema.size_z} Blocks / ${schema.total_size} bytes
        Description:
        \`\`\`
        ${schema.description}
        \`\`\`
        `
      };

      axios.post(schema_feed_url, data);
      //TODO: check status
    });
  }
};
