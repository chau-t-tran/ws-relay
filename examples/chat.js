socket = null;
room = "";

function displayMessage(message) {
	const messages = document.getElementById("messages");
	const newMessage = document.createElement("li");
	const time = new Date().toLocaleTimeString();
	const formattedMessage = `${time}: ${message}`;
	newMessage.appendChild(document.createTextNode(formattedMessage));
	messages.appendChild(newMessage);
}

function subscribeToRoom(key) {
	socket = new WebSocket(`ws://localhost:5000/${key}`);
	socket.addEventListener('message', (e) => {
		displayMessage(`RELAY: ${e.data}`);
	});
	displayMessage(`Connected to room ${key}`);
}

function sendMessage(text) {
	if (socket) {
		socket.send(text);
		displayMessage(`SELF: ${text}`);
	} else {
		displayMessage("Not connected to room! Join with \"/join sessionKey\"");
	}
}

function handleCommand(e) {
	const textInput = document.getElementById("message-input");
	const text = textInput.value;
	const tokens = text.split(" ");
	const op = tokens[0];

	if (op !== "/join") return sendMessage(text);
	if (tokens.length < 2) return;

	subscribeToRoom(tokens[1]);
}

window.addEventListener("load", function(event) {
	const button = document.getElementById("send-button");
	button.addEventListener("click", handleCommand);
});
