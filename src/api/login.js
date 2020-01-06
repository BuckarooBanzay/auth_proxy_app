
const app = require("../app");


const doLogin = require("../promise/login");

const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

app.post('/api/login', jsonParser, function(req, res){

	if (!req.body.username || !req.body.password){
		res.status(500).end();
		return;
	}

	doLogin(req.body.username, req.body.password)
	.then(result => {
		res.json({
			success: result.success,
			message: result.message,
		});
	})
	.catch(() => {
		res.status(500).end();
	});


});
