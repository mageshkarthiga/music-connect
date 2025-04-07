<template>
    <div class="chat-container">
        <!-- User List -->
        <div class="user-list">
            <h3>Users</h3>
            <ul>
                <li v-for="user in filteredUsers" :key="user.user_id" @click="joinRoom(user.user_id)"
                    :class="{ active: currentChatUser && currentChatUser.user_id === user.user_id }">
                    <img :src="user.profile_photo_url" alt="Profile" class="profile-photo" />
                    {{ user.user_name.charAt(0).toUpperCase() + user.user_name.slice(1) }}
                </li>
            </ul>
        </div>

        <!-- Chat Rooms -->
        <div class="chat-window" v-if="currentRoom && rooms.length > 0">
            <div v-for="room in rooms" :key="room.name" v-show="room.name === currentRoom" class="chat-room">
                <!-- Error Message -->
                <Message v-if="errorMessage" severity="error" :content="errorMessage" class="error-message" />
                <Card class="chat-card">
                    <template #header>
                        <div class="chat-header">
                            <h3>💬 Chat With: {{ room.name }}</h3>
                            <Button label="Leave Chat" class="p-button-danger" @click="leaveRoom(room)" />
                        </div>
                    </template>

                    <template #content>
                        <div class="chat-body" ref="chatBody">
                            <div v-for="(message, index) in room.messages" :key="index" class="message-wrapper"
                                :class="{ 'sent': message.isSent, 'received': !message.isSent }">
                                <div class="message-bubble">
                                    {{ message.message }} {{ message.timestamp }}
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

        <!-- Loading -->
        <div v-if="loading" class="loading-overlay">
        </div>
    </div>
</template>

<script>
export default {
    name: "ChatWindow",
    data() {
        return {
            currentUser: {
                user_id: 1,
                user_name: 'karthiga',
                profile_photo_url: 'https://ui-avatars.com/api/?name=Karthiga'
            },
            users: [
                {
                    user_id: 2,
                    user_name: 'alex',
                    profile_photo_url: 'https://ui-avatars.com/api/?name=Alex'
                },
                {
                    user_id: 3,
                    user_name: 'jordan',
                    profile_photo_url: 'https://ui-avatars.com/api/?name=Jordan'
                },
                {
                    user_id: 4,
                    user_name: 'charles',
                    profile_photo_url: 'https://ui-avatars.com/api/?name=Charles'
                }
            ],
            rooms: [],
            roomInput: "",
            ws: null,
            sockets: {},
            currentRoom: null,
            loading: false,
        };
    },
    computed: {
        filteredUsers() {
            return this.users.filter(user => user.user_id !== this.currentUser.user_id);
        }
    },
    mounted() {
        window.chatWindow = this // Expose the chat window instance globally for testing
    },
    methods: {
        connectToWebsocket() {
            this.ws = new WebSocket(`ws://localhost:8080/ws?userid=${this.currentUser.user_id}&room=${this.currentRoom}`);

            this.ws.onopen = () => {
                console.log("WebSocket connection established");
            };

            this.ws.onclose = () => {
                console.log("WebSocket connection closed");
                setTimeout(() => {
                    this.connectToWebsocket();
                }, 5000);
            };

            this.ws.onerror = (error) => {
                console.error("WebSocket error:", error);
                this.errorMessage = "WebSocket connection error. Please try again later.";
            };

            this.ws.addEventListener("message", (event) => {
                console.log("Message received:", event.data);
                this.handleNewMessage(event);
            });
        },

        handleNewMessage(event) {
            console.log("Raw WebSocket message:", event.data); // Log raw WebSocket message
            let data = event.data;
            data = data.split(/\r?\n/);

            for (let i = 0; i < data.length; i++) {
                try {
                    let msg = JSON.parse(data[i]);
                    console.log("Parsed message:", msg); // Log parsed message

                    // Validate the message structure
                    if (!msg.message || !msg.target || !msg.timestamp) {
                        console.warn("Invalid message received:", msg);
                        continue;
                    }

                    const room = this.findRoom(msg.target);
                    if (room && room.name === msg.target) {
                        console.log("Room found:", room.name); // Log the room being updated
                        this.$set(room, 'messages', [
                            ...room.messages,
                            {
                                message: msg.message.trim(),
                                sender: msg.sender || "Unknown",
                                isSent: msg.sender === this.currentUser.user_name,
                                timestamp: msg.timestamp
                            }
                        ]);
                        console.log(`Updated messages for room ${room.name}:`, room.messages); // Log updated messages
                    } else {
                        console.warn("Room not found for target:", msg.target);
                    }
                } catch (error) {
                    console.error("Error parsing message:", error);
                }
            }
        },

        sendMessage(room) {
            const ws = this.sockets[room.name];
            if (ws && ws.readyState === WebSocket.OPEN && room.newMessage.trim()) {
                const messageData = {
                    action: "send-message",
                    message: room.newMessage,
                    target: room.name,
                    timestamp: new Date().toISOString(),
                };

                ws.send(JSON.stringify(messageData));
                console.log("Message sent:", messageData);

                room.messages.push({
                    message: room.newMessage,
                    sender: this.currentUser.user_name,
                    isSent: true
                });

                room.newMessage = '';
            } else {
                console.error("WebSocket is not open. Cannot send message.");
                this.errorMessage = "Error: Cannot send message. Please try again later.";
            }
        },

        findRoom(roomName) {
            return this.rooms.find((room) => room.name === roomName);
        },

        joinRoom(userId) {
            const user = this.users.find(u => u.user_id === userId);

            // Check if the user exists
            if (!user) {
                console.error(`User with user_id ${userId} not found`);
                this.errorMessage = `User with user_id ${userId} not found`; // Set error message
                return;
            }

            const roomName = [this.currentUser.user_id, userId].sort().join("-");

            console.log(`User ${this.currentUser.user_name} is joining room: ${roomName}`);

            // Check if room already exists
            const existingRoom = this.rooms.find(r => r.name === roomName);
            if (existingRoom) {
                this.currentRoom = roomName;  // Set the current room when it's found
                this.currentChatUser = user;
                return;
            }

            const newRoom = {
                name: roomName,
                messages: [],
                newMessage: ''
            };

            this.loading = true;

            // Establish WebSocket connection for the room
            const ws = new WebSocket(`ws://localhost:8080/ws?userid=${this.currentUser.user_id}&room=${roomName}`);

            ws.onopen = () => {
                console.log(`WebSocket connection established for room: ${roomName}`);
                this.loading = false;

                const joinRoomMessage = {
                    action: "join-room",
                    message: roomName,
                    target: roomName
                };
                ws.send(JSON.stringify(joinRoomMessage));
                console.log("Join room message sent:", joinRoomMessage);
            };

            ws.onmessage = (event) => {
                const data = JSON.parse(event.data);
                newRoom.messages.push({
                    message: data.message,
                    sender: data.sender,
                    isSent: data.sender === this.currentUser.user_name
                });

                this.$nextTick(() => {
                    const body = this.$refs.chatBody;
                    if (body && body.scrollTop !== undefined) {
                        body.scrollTop = body.scrollHeight;
                    }
                });
            };

            ws.onclose = () => {
                console.log(`WebSocket connection closed for room: ${roomName}`);
                this.loading = false;
            };

            ws.onerror = (error) => {
                console.error(`WebSocket error for room: ${roomName}`, error);
                this.errorMessage = `Error connecting to room: ${roomName}`;
                this.loading = false;
            };

            this.rooms.push(newRoom);
            this.sockets[roomName] = ws;
            this.currentChatUser = user;
            this.currentRoom = roomName;
        },

        leaveRoom(room) {
            const socket = this.sockets[room.name];
            if (socket) {
                socket.close();
                delete this.sockets[room.name];
            }
            this.rooms = this.rooms.filter(r => r.name !== room.name);
            this.currentChatUser = null;
        },
    },
};
</script>

<style scoped>
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

.message-time {
    display: block;
    font-size: 0.75rem;
    color: #6c757d;
    margin-top: 0.25rem;
    text-align: right;
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

.loading-overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
}

.loading-content {
    width: 80%;
    max-width: 400px;
    background: white;
    padding: 1rem;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.error-message {
    margin: 1rem 0;
    text-align: center;
}

.message-sender {
    font-size: 0.8rem;
    color: #6c757d;
}
</style>