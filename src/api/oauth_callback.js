const app = require("../app");
const logger = require("../logger");
const axios = require('axios');

app.get('/api/oauth_callback', function(req, res){
  logger.debug("GET /api/oauth_callback", req.query);

  const data = {
    client_id: "68c2728e22f3a4b02dc0",
    client_secret: "xxx",
    code: req.query.code
  };
  axios.post("https://github.com/login/oauth/access_token", data)
  .then(r => console.log(r));


	res.end();
});
