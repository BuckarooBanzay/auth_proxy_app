

const events = require("../events");

module.exports = (username, password) => new Promise(function(resolve, reject){

	events.emit("channel-send", {
		type: "auth",
		data: {
			name: username,
			password: password
		}
	});

	function handleEvent(result){
		if (result.type == "auth" && result.data && result.data.name == username){
			events.removeListener("channel-recv", handleEvent);
			clearTimeout(handle);
			resolve(result.data);
		}
	}

	events.on("channel-recv", handleEvent);

	var handle = setTimeout(function(){
		events.removeListener("channel-recv", handleEvent);
		reject("mod-comm timeout");
	}, 1000);
});