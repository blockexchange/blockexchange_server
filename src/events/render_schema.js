const { spawn } = require('child_process');

const events = require("../events");

events.on("new-schema", schema => {
	console.log("Dispatching schema rendering");

	const child = spawn('node', ["src/worker/render_schema.js", schema.id]);

	child.stdout.on('data', (data) => {
	  console.log(`stdout: ${data}`);
	});

	child.stderr.on('data', (data) => {
	  console.error(`stderr: ${data}`);
	});

	child.on('close', (code) => {
  	console.log(`child process exited with code ${code}`);
		if (code == 0){
			events.emit("preview-rendered", schema.id);
		}
	});
});
