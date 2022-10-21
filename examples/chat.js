socket = null;

function subscribeToRoom(key) {
	socket = new WebSocket(`ws://localhost:5000/${key}`);
	socket.addEventListener('message', (e) => {
		console.log("FROM RELAY:", e.data);
	});
}

function sendMessage(text) {
	if (socket) {
		socket.send(text);
	}
}

function handleCommand(e) {
	const textInput = document.getElementById("message-input");
	const text = textInput.value;
	const tokens = text.split(" ");
	const op = tokens[0];

	if (op !== "/join") sendMessage(text);
	if (tokens.length < 2) return;

	subscribeToRoom(tokens[1]);
}

window.addEventListener("load", function(event) {
	const button = document.getElementById("send-button");
	button.addEventListener("click", handleCommand);
});
