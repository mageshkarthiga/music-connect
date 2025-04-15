<template>
    <div class="chat-container">
        <!-- User List -->
        <div class="user-list">
            <h3 class="user-list-header">
                <span>Users</span>
                <Button icon="pi pi-plus" rounded @click="redirectToUserTable()" />
            </h3>
            <ul>
                <li v-for="user in users" :key="user.user_id" @click="joinRoom(user)"
                    :class="{ active: currentChatUser && currentChatUser.user_id === user.user_id }">
                    <Avatar :image="user.profile_photo_url || '/public/profile.svg'" shape="circle" size="large" />
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
                            <h3>💬 Chat With: {{ room.otherUserName.charAt(0).toUpperCase() +
                                room.otherUserName.slice(1) || "Loading..." }}</h3>
                        </div>
                    </template>

                    <template #content>
                        <div v-if="loading" class="card-spinner-container">
                            <i class="pi pi-spin pi-spinner card-spinner"></i>
                        </div>
                        <div class="chat-body" ref="chatBody">
                            <div v-if="room.messages.length === 0" class="no-messages">
                                It's quiet here… start the conversation and share the vibes 🎧✨
                            </div>
                            <!-- Show messages if they exist -->
                            <div v-else>
                                <div v-for="(message, index) in room.messages" :key="index" class="message-wrapper"
                                    :class="{ 'sent': message.isSent, 'received': !message.isSent }">
                                    <div class="message-bubble">
                                        {{ message.message }}
                                    </div>
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
import UserService from "@/service/UserService";
import axios from "axios";

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
            users: [],
            rooms: [],
            roomInput: "",
            ws: null,
            sockets: {},
            currentRoom: null,
            loading: false,
        };
    },
    watch: {
        selectedUserId: {
            immediate: true,
            async handler(newUserId) {
                if (newUserId) {
                    const user = this.users.find((user) => user.firebase_uid === newUserId);

                    if (user) {
                        await this.joinRoom(user);
                    } else {
                        try {
                            const response = await axios.get(`http://localhost:8080/firebase/${newUserId}`, {
                                withCredentials: true,
                            });
                            const fetchedUser = response.data;

                            this.users.push(fetchedUser);

                            await this.joinRoom(fetchedUser);
                        } catch (error) {
                            console.error(`Error fetching user with ID ${newUserId}:`, error);
                            this.errorMessage = "Failed to fetch user details. Please try again.";
                        }
                    }
                }
            },
        },
    },
    mounted() {
        console.log("Selected User ID:", this.selectedUserId);
        this.getCurrentUser().then(async () => {
            await this.fetchChatHistoryUsers();

            let user = this.users.find(user => user.firebase_uid === this.selectedUserId);

            if (!user && this.selectedUserId) {
                try {
                    const response = await UserService.getUserByFirebaseUID(this.selectedUserId);
                    user = response.data;
                    this.users.push(user);
                } catch (error) {
                    console.warn(`User with ID ${this.selectedUserId} could not be fetched.`, error);
                    return;
                }
            }

            if (user) {
                this.joinRoom(user);
            } else {
                console.warn("No user selected or user not found.");
            }
        });
    },
    methods: {
        async getCurrentUser() {
            try {
                const response = await UserService.getUser();
                this.currentUser.user_id = response.firebase_uid;
                this.currentUser.user_name = response.userName;
            } catch (error) {
                console.error("Error fetching current user:", error);
                this.errorMessage = "Failed to fetch current user. Please try again.";
            }
        },
        async getOtherUsers(userID) {
            try {
                const response = await UserService.getUserByFirebaseUID(userID)
                console.log(response)
                return response.userName;
            }
            catch(error) {
                console.error("Error getting user:", error);
                return "Unknown User";
            };
        },
        async fetchChatHistoryUsers() {
            try {
                const response = await axios.get(`http://localhost:8080/users/${this.currentUser.user_id}/chat-history`, {
                    withCredentials: true,
                });
                const userIds = response.data;

                const userDetailsPromises = userIds.map(async (userId) => {
                    try {
                        const userResponse = await axios.get(`http://localhost:8080/firebase/${userId}`, {
                            withCredentials: true,
                        });
                        return userResponse.data;
                    } catch (error) {
                        console.error(`Error fetching details for user ID ${userId}:`, error);
                        return null;
                    }
                });

                const userDetails = await Promise.all(userDetailsPromises);

                this.users = userDetails.filter((user) => user !== null);
            } catch (error) {
                console.error("Error fetching users with chat history:", error);
                this.errorMessage = "Failed to load chat history users.";
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
                            this.$nextTick(() => {
                                const chatBody = this.$refs.chatBody;
                                if (chatBody) {
                                    chatBody.scrollTop = chatBody.scrollHeight;
                                }
                            });
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

                room.newMessage = '';
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
            console.log(`User ${this.currentUser.user_name} is joining room: ${roomName}`);

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
                otherUserName: "Loading...", // Placeholder for the other user's name
            };

            this.loading = true;

            try {
                // Fetch the other user's name
                const otherUserID = roomName.split("-").filter(id => id !== this.currentUser.user_id)[0];
                const response = await this.getOtherUsers(otherUserID);
                newRoom.otherUserName = response;

                // Fetch chat history via REST API
                const messagesResponse = await axios.get(`http://localhost:8080/rooms/${roomName}/messages`,
                    { withCredentials: true }
                );
                newRoom.messages = messagesResponse.data.map(msg => ({
                    message: msg.message,
                    sender: msg.sender,
                    isSent: msg.sender === this.currentUser.user_id,
                }));
                console.log(`Fetched messages for room ${roomName}:`, newRoom.messages);
            } catch (error) {
                console.error(`Error fetching messages for room ${roomName}:`, error);
                this.errorMessage = "Failed to load chat history. Please try again later.";
            }

            // Establish WebSocket connection for the room
            const ws = new WebSocket(`ws://localhost:8080/ws?userid=${this.currentUser.user_id}&room=${roomName}`);

            ws.onopen = () => {
                console.log(`WebSocket connection established for room: ${roomName}`);
                this.loading = false;

                const joinRoomMessage = {
                    action: "join-room",
                    message: roomName,
                    target: roomName,
                };
                ws.send(JSON.stringify(joinRoomMessage));
                console.log("Join room message sent:", joinRoomMessage);
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

            this.rooms.push(newRoom);
            this.sockets[roomName] = ws;
            this.currentChatUser = user;
            this.currentRoom = roomName;
        },
        redirectToUserTable() {
            this.$router.push({ name: "search" });
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
    max-width: 300px;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1rem;
    background-color: #ffffff;
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
    background-color: #e0e0e0;
}


.user-list li.active {
    background-color: #d1e7dd;
}

.user-list h3 {
    padding: 0.5rem 0;
    margin: 0;
    font-size: 1.2rem;
    border-bottom: 1px solid #ddd;
    background-color: #ffffff;
    position: static;
    top: auto;
    z-index: auto;
}

.chat-window {
    flex: 2;
}

.chat-card {
    margin-top: 50px;
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

.user-list-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
    margin: 0;
    font-size: 1.2rem;
    border-bottom: 1px solid #ddd;
    background-color: #ffffff;
    position: static;
    top: auto;
    z-index: auto;
}
</style>