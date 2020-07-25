const events = require("../events");
const user_dao = require("../dao/user");
const axios = require('axios');

module.exports = function(){
  // listen to events and dispatch webhook

  const schema_feed_url = process.env.DISCORD_SCHEMA_FEED_URL;

  if (schema_feed_url){
    events.on("new-schema", function(schema){
      //https://birdie0.github.io/discord-webhooks-guide/examples/spotify.html
      user_dao.get_by_id(schema.user_id)
      .then(function(user){
        const data = {
          content: `Schema created: **${schema.name}** by **${user.name}**\n` +
            `License: **${schema.license}**\n` +
            `Size: ${schema.size_x}/${schema.size_y}/${schema.size_z} Blocks / ${schema.total_size} bytes\n` +
            `Description:\n\`\`\`\n${schema.description}\n\`\`\`\n`
        };

        axios.post(schema_feed_url, data);
        //TODO: check status
      });
    });
  }
};
