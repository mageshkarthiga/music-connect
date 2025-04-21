<template>
    <div class="chat-container">
        <!-- Friend List -->
        <div class="user-list">
            <h3 class="list-title">Friends</h3>
            <ul v-if="friends.length > 0">
                <li v-for="friend in friends" :key="friend.user_id" @click="joinRoom(friend)"
                    :class="{ active: currentChatUser && currentChatUser.user_id === friend.user_id }">
                    <Avatar :image="friend.profile_photo_url || '/profile.svg'" shape="circle" size="large" />
                    {{ friend.user_name.charAt(0).toUpperCase() + friend.user_name.slice(1) }}
                </li>
            </ul>
            <div v-else>
                <p>No friends available. Connect with friends to start chatting!✨</p>
            </div>
        </div>

        <!-- Chat Rooms -->
        <div class="chat-window" v-if="currentRoom && rooms.length > 0">
            <div v-for="room in rooms" :key="room.name" v-show="room.name === currentRoom" class="chat-room">
                <!-- Error Message -->
                <Message v-if="errorMessage" severity="error" :content="errorMessage" class="error-message" />
                <Card class="chat-card">
                    <template #header>
                        <div class="chat-header">
                            <Avatar :image="getFriendById(getOtherUserId(room.name))?.profile_photo_url || '/profile.svg'" shape="circle"
                                size="large" class="chat-header-avatar" />
                            <h3 class="chat-header-name">
                                {{ room.otherUserName.charAt(0).toUpperCase() +
                                    room.otherUserName.slice(1) || "Loading..." }}
                            </h3>
                        </div>
                    </template>

                    <template #content>
                        <div v-if="loading" class="card-spinner-container">
                            <i class="pi pi-spin pi-spinner card-spinner"></i>
                        </div>

                        <div class="chat-body-wrapper">
                            <div class="chat-body" :ref="el => setChatBodyRef(room.name, el)">
                                <div v-if="room.messages.length === 0" class="no-messages">
                                    It's quiet here… start the conversation and share the vibes 🎧✨
                                </div>
                                <div v-else>
                                    <div v-for="(message, index) in room.messages" :key="index" class="message-wrapper"
                                        :class="{ 'sent': message.isSent, 'received': !message.isSent }">
                                        <div class="message-bubble">
                                            {{ message.message }}
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <Button v-if="room.showScrollArrow" icon="pi pi-arrow-down" class="scroll-to-bottom"
                                @click="scrollToBottom(room.name)" />
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
import UserService from "@/service/UserService";
import axios from "axios";
import FriendService from "@/service/FriendService";

export default {
    name: "ChatWindow",
    props: {
        selectedUserId: {
            type: String,
            required: false,
        },
    },
    data() {
        return {
            currentUser: {
                user_id: null,
                user_name: null,
            },
            friends: [],
            rooms: [],
            roomInput: "",
            ws: null,
            sockets: {},
            currentRoom: null,
            loading: false,
            chatBodies: {},
            errorMessage: "",
            CHAT_URL: process.env.VUE_APP_CHAT_URL,
        };
    },
    async mounted() {
        await this.getCurrentUser();
        await this.fetchFriends();

        if (this.selectedUserId) {
            const user = this.friends.find(user => user.firebase_uid === this.selectedUserId);

            if (user) {
                this.joinRoom(user);
            } else {
                console.warn(`User with ID ${this.selectedUserId} not found in chat history.`);

                try {
                    const response = await UserService.getUserByFirebaseUID(this.selectedUserId);
                    if (response) {
                        this.friends.push(response);
                        this.joinRoom(response);
                    } else {
                        console.warn(`Unable to fetch details for user ID ${this.selectedUserId}`);
                    }
                } catch (error) {
                    console.error("Error fetching user details:", error);
                }
            }
        } else {
            console.log("No user selected. Displaying user list.");
        }
    },
    methods: {
        async getCurrentUser() {
            try {
                const response = await UserService.getUser();
                this.currentUser.user_id = response.firebase_uid;
                this.currentUser.user_name = response.user_name || "Anonymous";
            } catch (error) {
                console.error("Error getting user:", error);
            }
        },
        async getOtherUsers(userID) {
            try {
                const response = await UserService.getUserByFirebaseUID(userID);
                return response.user_name
            } catch (error) {
                console.error("Error getting user:", error);
                return "Unknown User";
            }
        },
        async fetchFriends() {
            try {
                const response = await FriendService.getFriends();
                this.friends = response;
                console.log("Fetched friends:", this.friends);
            } catch (error) {
                console.error("Error fetching friends:", error);
                this.errorMessage = "Failed to load friends. Please try again later.";
            }
        },
        handleNewMessage(event) {
            console.log("Raw WebSocket message:", event.data);
            let data = event.data;
            data = data.split(/\r?\n/);

            for (let i = 0; i < data.length; i++) {
                try {
                    let msg = JSON.parse(data[i]);
                    console.log("Parsed message:", msg);

                    // Validate the message structure
                    if (!msg.message || !msg.target) {
                        console.warn("Invalid message received:", msg);
                        continue;
                    }

                    const room = this.findRoom(msg.target);
                    if (room && room.name === msg.target) {
                        // Check if we already have this message (by message_id)
                        const isDuplicate = room.messages.some(existingMsg =>
                            msg.message_id && existingMsg.message_id === msg.message_id
                        );

                        if (!isDuplicate) {
                            console.log("Adding new message to room:", room.name);
                            room.messages.push({
                                message: msg.message.trim(),
                                sender: msg.sender || "Unknown",
                                isSent: msg.sender === this.currentUser.user_id,
                                message_id: msg.message_id,
                            });

                            // Scroll to the bottom of the chat
                            if (!this.isChatScrolledToBottom(room.name)) {
                                room.showScrollArrow = true;
                            } else {
                                this.$nextTick(() => {
                                    this.scrollToBottom(room.name);
                                });
                            }

                        } else {
                            console.log("Skipping duplicate message:", msg.message_id);
                        }
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
                };

                ws.send(JSON.stringify(messageData));
                console.log("Message sent:", messageData);

                room.messages.push({
                    message: room.newMessage,
                    sender: this.currentUser.user_name,
                    isSent: true,
                });

                room.newMessage = '';  // Clear input
            } else {
                console.error("WebSocket is not open. Cannot send message.");
                this.errorMessage = "Error: Cannot send message. Please try again later.";
            }
        },


        findRoom(roomName) {
            return this.rooms.find((room) => room.name === roomName);
        },

        async joinRoom(user) {
            const roomName = [this.currentUser.user_id, user.firebase_uid].sort().join("-");
            console.log("Joining room for user:", user);
            console.log("Generated room name:", roomName);

            const existingRoom = this.rooms.find(r => r.name === roomName);
            if (existingRoom) {
                this.currentRoom = roomName;
                this.currentChatUser = user;
                return;
            }

            const newRoom = {
                name: roomName,
                messages: [],
                newMessage: '',
                otherUserName: user.user_name || "Loading...",
                showScrollArrow: false,
            };

            this.loading = true;

            try {
                const otherUserID = roomName.split("-").filter(id => id !== this.currentUser.user_id)[0];
                const response = await this.getOtherUsers(otherUserID);
                newRoom.otherUserName = response;

                const messagesResponse = await axios.get(`${this.CHAT_URL}/rooms/${roomName}/messages`, { withCredentials: true });
                if (messagesResponse.data.length > 0) {
                    newRoom.messages = messagesResponse.data.map(msg => ({
                        message: msg.message,
                        sender: msg.sender,
                        isSent: msg.sender === this.currentUser.user_id,
                    }));
                    console.log(`Fetched messages for room ${roomName}:`, newRoom.messages);
                } else {
                    console.log(`No previous messages for room ${roomName}`);
                }
            } catch (error) {
                console.error(`Error fetching messages for room ${roomName}:`, error);
                this.errorMessage = "Failed to load chat history. Please try again later.";
            }

            const ws = new WebSocket(this.getWebSocketURL(this.currentUser.user_id, roomName));

            ws.onopen = () => {
                console.log(`WebSocket connection established for room: ${roomName}`);
                this.loading = false;

                const joinRoomMessage = {
                    action: "join-room",
                    message: roomName,
                    target: roomName,
                };
                ws.send(JSON.stringify(joinRoomMessage));
            };

            ws.onmessage = (event) => {
                this.handleNewMessage(event);
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

            // Push the new room to the rooms array
            this.rooms.push(newRoom);
            this.sockets[roomName] = ws;

            // Set the current room and the chat user
            this.currentChatUser = user;
            this.currentRoom = roomName;

            this.$nextTick(() => {
                this.scrollToBottom(roomName);
            });
        },
        setChatBodyRef(roomName, el) {
            if (el) {
                this.chatBodies[roomName] = el;

                el.addEventListener('scroll', () => {
                    const nearBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 30;
                    const room = this.findRoom(roomName);
                    if (room) {
                        room.showScrollArrow = !nearBottom;
                    }
                });
            }
        },
        scrollToBottom(roomName) {
            this.$nextTick(() => {
                const chatBody = this.chatBodies[roomName];
                if (chatBody) {
                    chatBody.scrollTop = chatBody.scrollHeight;
                }
            });
        },
        isChatScrolledToBottom(roomName) {
            const el = this.chatBodies[roomName];
            return el && el.scrollHeight - el.scrollTop - el.clientHeight < 30;
        },
        getWebSocketURL(userId, roomName) {
            const isLocalhost = false;
            const protocol = isLocalhost ? 'ws' : 'wss';
            const host = isLocalhost ? 'localhost:8080' : this.CHAT_URL.replace(/^https?:\/\//, '');

            return `${protocol}://${host}/ws?userid=${userId}&room=${roomName}`;
        },
        getOtherUserId(roomName) {
            const ids = roomName.split('-');
            return ids.find(id => id !== this.currentUser.user_id);
        },
        getFriendById(userId) {
            return this.friends.find(friend => friend.firebase_uid === userId);
        }
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
    max-width: 300px;
    border: 1px solid var(--border-light);
    border-radius: 8px;
    padding: 1rem;
    background-color: var(--primary-bg-light);
    max-height: 600px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    margin-top: 50px;
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
    background-color: rgba(0, 0, 0, 0.05);
}


.user-list li.active {
    background-color: #d1e7dd;
}

.user-list h3 {
    padding: 0.5rem 0;
    margin: 0;
    font-size: 1.2rem;
    border-bottom: 1px solid var(--border-light);
    background-color: var(--primary-bg-light);
    position: static;
    top: auto;
    z-index: auto;
}

.chat-window {
    flex: 2;
}

.chat-card {
    margin-top: 60px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chat-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.5rem 1rem;
    background-color: #f5f5f5;
    border-radius: 8px 8px 0 0;
}

.chat-header-avatar {
    flex-shrink: 0;
}

.chat-header-name {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    color: #343a40;
}

.chat-body {
    max-height: 400px;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    padding: 1rem;
}

.separator {
    border: none;
    border-bottom: 1px solid var(--border-light);
    margin: 0 0 1rem 0;
}

.message-wrapper {
    display: flex;
    margin-bottom: 0.5rem;
}

.message-wrapper.sent {
    justify-content: flex-end;
}

.message-wrapper.received {
    justify-content: flex-start;
}

.message-bubble {
    padding: 0.75rem 1.25rem;
    border-radius: 8px;
    max-width: 70%;
    word-break: break-word;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
    font-size: 0.95rem;
    line-height: 1.4;
}

.message-wrapper.sent .message-bubble {
    background-color: #d1e7dd;
    color: #0f5132;
    border: 1px solid #badbcc;
}

.message-wrapper.received .message-bubble {
    background-color: #ffffff;
    color: #343a40;
    border: 1px solid #dee2e6;
}

.message-sender {
    font-size: 0.75rem;
    color: #6c757d;
    margin-bottom: 0.25rem;
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

.card-spinner-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.card-spinner {
    font-size: 2rem;
    color: #6c757d;
}

.error-message {
    margin: 1rem 0;
    text-align: center;
}

.message-sender {
    font-size: 0.8rem;
    color: #6c757d;
}

.no-messages {
    text-align: center;
    font-size: 1rem;
    color: #6c757d;
    font-style: italic;
    margin-top: 1rem;
}

.chat-body-wrapper {
    position: relative;
}

.scroll-to-bottom {
    position: absolute;
    bottom: 16px;
    left: 50%;
    transform: translateX(-50%);
    z-index: 10;
    color: white;
    border-radius: 50%;
    width: 2.5rem;
    height: 2.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}
</style>