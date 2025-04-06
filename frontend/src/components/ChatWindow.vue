<template>
    <div class="chat-container">
        <!-- User List -->
        <div class="user-list">
            <h3>Users</h3>
            <ul>
                <li v-for="user in users" :key="user.user_id" @click="joinRoom(user.user_name)"
                    :class="{ active: currentChatUser && currentChatUser.user_id === user.user_id }">
                    <img :src="user.profile_photo_url" alt="Profile" class="profile-photo" />
                    {{ user.user_name }}
                </li>
            </ul>
        </div>

        <!-- Chat Rooms -->
        <div class="chat-window" v-if="rooms.length > 0">
            <div v-for="room in rooms" :key="room.name" class="chat-room">
                <Card class="chat-card">
                    <template #header>
                        <div class="chat-header">
                            <h3>💬 Chat Room: {{ room.name }}</h3>
                            <Button label="Leave Room" class="p-button-danger" @click="leaveRoom(room)" />
                        </div>
                    </template>

                    <template #content>
                        <div class="chat-body" ref="chatBody">
                            <div v-for="(message, index) in room.messages" :key="index" class="message-wrapper"
                                :class="{ 'sent': message.isSent, 'received': !message.isSent }">
                                <div class="message-bubble">
                                    {{ message.message }}
                                    <span v-if="message.sender" class="msg-sender">{{ message.sender }}</span>
                                </div>
                            </div>
                        </div>
                    </template>

                    <template #footer>
                        <div class="chat-footer">
                            <Textarea v-model="room.newMessage" rows="2" autoResize placeholder="Type your message..."
                                @keyup.enter.exact="sendMessage(room)" class="chat-input" />
                            <Button icon="pi pi-send" class="p-button-rounded p-button-primary"
                                @click="sendMessage(room)" aria-label="Send" />
                        </div>
                    </template>
                </Card>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "ChatWindow",
    data() {
        return {
            currentUser: {
                user_id: 4, // Simulated current user ID
                user_name: "You", // Simulated current user name
                profile_photo_url: "https://via.placeholder.com/50", // Simulated profile photo
            },
            users: [
                {
                    user_id: 1,
                    user_name: "Alice",
                    profile_photo_url: "https://via.placeholder.com/50",
                },
                {
                    user_id: 2,
                    user_name: "Bob",
                    profile_photo_url: "https://via.placeholder.com/50",
                },
                {
                    user_id: 3,
                    user_name: "Charlie",
                    profile_photo_url: "https://via.placeholder.com/50",
                },
            ],
            rooms: [], // List of chat rooms
            roomInput: "", // Input for joining a room
            ws: null, // WebSocket connection
        };
    },
    mounted() {
        this.connectToWebsocket();
    },
    methods: {
        connectToWebsocket() {
            this.ws = new WebSocket("ws://localhost:8080/ws");

            this.ws.onopen = () => {
                console.log("WebSocket connection established");
            };

            this.ws.onclose = () => {
                console.log("WebSocket connection closed");
            };

            this.ws.onerror = (error) => {
                console.error("WebSocket error:", error);
            };

            this.ws.addEventListener("message", (event) => {
                console.log("Message received:", event.data); // Log incoming messages
                this.handleNewMessage(event);
            });
        },

        handleNewMessage(event) {
            let data = event.data;
            data = data.split(/\r?\n/);

            for (let i = 0; i < data.length; i++) {
                let msg = JSON.parse(data[i]);
                const room = this.findRoom(msg.target);
                if (typeof room !== "undefined") {
                    room.messages.push(msg);
                }
            }
        },

        sendMessage(room) {
            if (room.newMessage !== "") {
                const message = {
                    action: "send-message",
                    message: room.newMessage.trim(), // Trim the message to remove unnecessary spaces or newlines
                    target: room.name,
                };

                console.log("Sending message:", message); // Log the message being sent

                this.ws.send(JSON.stringify(message));
                room.messages.push({
                    message: room.newMessage,
                    isSent: true,
                    sender: "You",
                });

                room.newMessage = "";
            }
        },

        findRoom(roomName) {
            return this.rooms.find((room) => room.name === roomName);
        },

        joinRoom(roomName) {
            this.ws.send(JSON.stringify({ action: "join-room", message: roomName }));
            this.rooms.push({ name: roomName, messages: [], newMessage: "" });
        },

        leaveRoom(room) {
            this.ws.send(JSON.stringify({ action: "leave-room", message: room.name }));
            this.rooms = this.rooms.filter((r) => r.name !== room.name);
        },
    },
};
</script>

<style scoped>
/* Add your styles here */
.chat-container {
    display: flex;
    max-width: 900px;
    margin: 0 auto;
    gap: 1rem;
    padding: 1rem;
}

.user-list {
    flex: 1;
    max-width: 200px;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1rem;
    background-color: #f9f9f9;
}

.user-list ul {
    list-style: none;
    padding: 0;
    margin: 0;
}

.user-list li {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem;
    cursor: pointer;
    border-radius: 4px;
    transition: background-color 0.2s;
}

.user-list li:hover {
    background-color: #e9ecef;
}

.user-list li.active {
    background-color: #d1e7dd;
}

.profile-photo {
    width: 40px;
    height: 40px;
    border-radius: 50%;
}

.chat-window {
    flex: 2;
}

.chat-card {
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chat-header {
    text-align: center;
    font-size: 1.5rem;
    font-weight: 600;
    padding: 0.5rem 0;
    background-color: #f5f5f5;
    border-radius: 8px 8px 0 0;
}

.chat-body {
    max-height: 400px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding: 1rem;
}

.message-wrapper {
    display: flex;
}

.message-wrapper.sent {
    justify-content: flex-end;
}

.message-wrapper.received {
    justify-content: flex-start;
}

.message-bubble {
    padding: 0.6rem 1rem;
    border-radius: 1.25rem;
    max-width: 70%;
    word-break: break-word;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.message-wrapper.sent .message-bubble {
    background-color: #d1e7dd;
    color: #0f5132;
    border-bottom-right-radius: 0;
}

.message-wrapper.received .message-bubble {
    background-color: #f8d7da;
    color: #842029;
    border-bottom-left-radius: 0;
}

.chat-footer {
    display: flex;
    align-items: flex-end;
    gap: 0.5rem;
    padding: 1rem;
    background-color: #f5f5f5;
    border-radius: 0 0 8px 8px;
}

.chat-input {
    flex: 1;
}
</style>