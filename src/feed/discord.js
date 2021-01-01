const axios = require('axios');

const events = require("../events");

const user_dao = require("../dao/user");
const schema_dao = require("../dao/schema");
const schema_screenshot_dao = require("../dao/schema_screenshot");

const BASE_URL = process.env.BASE_URL;

module.exports = function(){
  // listen to events and dispatch webhook

  const schema_feed_url = process.env.DISCORD_SCHEMA_FEED_URL;

	if (!schema_feed_url){
		// no discord feed configured
		return;
	}

  events.on("preview-rendered", async function(schema_id){
		const schema = await schema_dao.get_by_id(schema_id);
		const previews = await schema_screenshot_dao.find_all(schema_id);

		let preview_txt = "";
		if (BASE_URL && previews && previews.length > 0){
			// show a preview, if available
			preview_txt = `Preview: ${BASE_URL}/api/schema/${schema_id}/screenshot/${previews[0].id}`;
		}

    //https://birdie0.github.io/discord-webhooks-guide/examples/spotify.html
    const user = await user_dao.get_by_id(schema.user_id);
    const data = {
      content: `Schema created: **${schema.name}** by **${user.name}**\n` +
				`Link: ${BASE_URL}/api/static/schema/${user.name}/${schema.name}\n` +
        `License: **${schema.license}**\n` +
        `Size: ${schema.max_x+1}/${schema.max_y+1}/${schema.max_z+1} Blocks / ${schema.total_size} bytes\n` +
				`Description:\n\`\`\`\n${schema.description}\n\`\`\`\n` +
				`Download:\n\`\`\`\n/bx_load ${user.name} ${schema.name}\n\`\`\`\n` +
				preview_txt
    };

    axios.post(schema_feed_url, data);
    //TODO: check status
  });
};
