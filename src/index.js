
const app = require("./app");
require("./api/channel");
require("./api/login");

app.listen(8080, () => console.log('Listening on http://127.0.0.1:8080'))
