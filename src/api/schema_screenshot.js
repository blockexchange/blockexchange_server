const bodyParser = require('body-parser');
var Canvas = require('canvas');

const jsonParser = bodyParser.json();
const logger = require("../logger");

const app = require("../app");
const schema_screenshot_dao = require("../dao/schema_screenshot");
const schema_dao = require("../dao/schema");

const tokenmiddleware = require("../middleware/token");
const permission_create = tokenmiddleware(claims => claims.permissions.screenshot.create);
const permission_delete = tokenmiddleware(claims => claims.permissions.screenshot.delete);

app.get('/api/schema/:id/screenshot', function(req, res){
  logger.debug("GET /api/schema/:id/screenshot", req.params.id);

  schema_screenshot_dao.find_all(req.params.id)
  .then(screenshots => screenshots || [])
  .then(screenshots => screenshots.map(s => ({ id: s.id, type: s.type, title: s.title })))
  .then(screenshots => res.json(screenshots))
  .catch(() => res.status(500).end());
});

app.get('/api/schema/:id/screenshot/:screenshot_id', function(req, res){
  logger.debug("GET /api/schema/:id/screenshot/:screenshot_id", req.params.id, req.params.screenshot_id);

  schema_screenshot_dao.get_by_id(+req.params.id, +req.params.screenshot_id)
  .then(screenshot => {
		if (req.query.width && req.query.height){
			// resize
			var img = new Canvas.Image();
			img.onerror = function(e){
				console.error(e);
			};
			img.onload = function(){
				const width = +req.query.width;
				const height = +req.query.height;
				//console.log(img.height);

				const canvas = Canvas.createCanvas(width, height);
				var ctx = canvas.getContext('2d');
 				ctx.drawImage(img, 0, 0, width, height);

				res.header("Content-type", "image/png")
		    .send(canvas.toBuffer("image/png"));
			};

			img.src = "data:image/png;base64," + screenshot.data.toString("base64");
		} else {
			// original size
			res.header("Content-type", screenshot.type)
	    .send(screenshot.data);
		}

  })
  .catch(e => {
		console.error(e);
		res.status(500).end();
	});
});

app.post('/api/schema/:id/screenshot', permission_create, jsonParser, function(req, res){
  logger.debug("POST /api/schema/:id/screenshot", req.params.id);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

    schema_screenshot_dao.create(schema.id, req.body.title, req.body.type, req.body.data)
    .then(() => res.end());
  });
});


app.get('/api/schema/:id/screenshot/:screenshot_id', permission_delete, function(req, res){
  logger.debug("DELETE /api/schema/:id/screenshot/:screenshot_id", req.params.id, req.params.screenshot_id);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

    schema_screenshot_dao.remove(req.params.screenshot_id)
    .then(() => res.end());
  })
  .catch(() => res.status(500).end());
});
