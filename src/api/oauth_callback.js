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

  const options = {
    headers: {
      "Accept": "application/json"
    }
  };

  axios.post("https://github.com/login/oauth/access_token", data, options)
  .then(r => {
    if (!r.data.access_token){
      console.log(r);
      throw new Error("No access_token received!");
    }
    return axios.get("https://api.github.com/user", {
      headers: {
        "Authorization": "Bearer " + r.data.access_token
      }
    });
  })
  .then(r => {
    const user_info = r.data;
    console.log(user_info);
    // user_info.login / avatar_url / name / email
    //TODO: build jwt
    res.redirect("http://127.0.0.1:8080/#/oauth?jwt=blah");
  })
  .catch(e => {
    console.log(e);
    res.send(e.messages);
  });
});
