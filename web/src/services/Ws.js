import Chats from "../views/ChatsView.js"
import Utils from "./Utils.js";
import Fetcher from "./Fetcher.js";

let connection

// Managing a WebSocket connection.
// Only one active connection is maintained.
const getConnection = () => {
/*
The WebSocket.readyState read-only property returns the current state of the WebSocket connection.
CONNECTING (0)
OPEN (1)
CLOSING (2)
CLOSED (3)
 */
    if (connection && connection.readyState < 2) {
        return Promise.resolve(connection)
    }

    return new Promise((resolve, reject) => {
        // Verifies if the browser supports WebSocket.
        if (window["WebSocket"]) {
            let token = localStorage.getItem("accessToken");
            console.log("The accessToken is: " + token);

            /*
            The token sent by the API server, and stored in the browser localstore
            using th signIn() function.
             */
            if (token === undefined) {
                alert("error opening websocket connection, no access token in localStorage")
                return
            }

            // Establishes a new WebSocket connection to the server. API_HOST_NAME defined in the index.html
            const conn = new WebSocket(`ws://${API_HOST_NAME}/ws`)

            // Send the access token to the server for authentication as soon as the connection is open.
            conn.onopen = function () {
                conn.send(JSON.stringify({ type: "token", body: token }))
            }

            // The connection error handling
            conn.onerror = function (evt) {
                Utils.showError(503)
                return
            }

            conn.onmessage = async function (evt) {
                let obj = JSON.parse(evt.data)
/*
Message Handling:

message: Draws a new message in the chat.
messagesResponse: Prepends messages to the chat.
chatsResponse: Draws the chat list.
readMessageResponse: Marks messages as read.
notification: Logs notifications.
onlineUsersResponse: Draws online users.
typingInResponse: Displays a typing indicator.
error: Handles errors like token expiration, attempting to refresh the token and reauthenticate.
successConnection: Resolves the promise, signaling that the connection was successful.
pingMessage: Responds to server pings with a pong message.
 */
                switch (obj.type) {
                    case "message":
                        await Chats.drawNewMessage(obj.body)
                        break
                    case "messagesResponse":
                        await Chats.prependMessages(obj.body)
                        break
                    case "chatsResponse":
                        Chats.drawChats(obj.body)
                        break
                    case "readMessageResponse":
                        await Chats.changeMessageStatusToRead(obj.body)
                        break
                    case "notification":
                        console.log(obj)
                        break
                    case "onlineUsersResponse":
                        await Chats.drawOnlineUsers(obj.body)
                        break
                    case "typingInResponse":
                        await Chats.drawTypingIn(obj.body)
                        break
                    case "error":
                        if (obj.body == "token has expired") {
                            await Fetcher.refreshToken()
                            token = localStorage.getItem("accessToken")
                            conn.send(JSON.stringify({ type: "token", body: token }))
                        } else {
                            alert(obj.body)
                        }
                        break
                    case "successConnection":
                        resolve(conn)
                        break
                    case "pingMessage":
                        conn.send(JSON.stringify({ type: "pongMessage" }))
                        break
                    default:
                        console.log(obj)
                        break
                }
            };
        } else {
            alert("Your browser does not support WebSockets")
        }
    })
}


const Ws = {
    connect: async () => {
        connection = await getConnection()
    },

    send: async (e) => {
        connection = await getConnection()
        connection.send(e)
    },

    disconnect: async () => {
        connection.close()
    }
}

export default Ws

