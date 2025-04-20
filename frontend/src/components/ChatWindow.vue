<template>
    <div class="chat-container">
        <!-- Friend List -->
        <div class="user-list">
            <h3>Friends</h3>
            <ul>
                <li v-for="friend in friends" :key="friend.user_id" @click="joinRoom(friend)"
                    :class="{ active: currentChatUser && currentChatUser.user_id === friend.user_id }">
                    <Avatar :image="friend.profile_photo_url || '/profile.svg'" shape="circle" size="large" />
                    {{ friend.user_name.charAt(0).toUpperCase() + friend.user_name.slice(1) }}
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
                        <h3 class="chat-header">
                            <div class="user-info">
                                <img class="user-avatar" :src="room.otherUserProfilePic || 'default-avatar.jpg'" alt="User's Avatar" />
                                <span class="user-name">{{ room.otherUserName.charAt(0).toUpperCase() + room.otherUserName.slice(1) || "Loading..." }}</span>
                            </div>
                        </h3>
                    </template>

                    <template #content>
                        <!-- Chat Content Goes Here -->
                    </template>

                    <template #footer>
                        <div class="chat-footer">
                            <Textarea v-model="room.newMessage" rows="2" autoResize placeholder="Type your message..." @keyup.enter.exact="sendMessage(room)" class="chat-input" />
                            <Button icon="pi pi-send" class="p-button-rounded p-button-primary" @click="sendMessage(room)" aria-label="Send" />
                        </div>
                    </template>
                </Card>
            </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="loading-overlay"></div>
    </div>
</template>


<script>
import UserService from "@/service/UserService";
import friendService from "@/service/FriendService";
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
            friends:[],
            users: [],
            rooms: [],
            roomInput: "",
            ws: null,
            sockets: {},
            currentRoom: null,
            loading: false,
            chatBodies: {},
            errorMessage: ""
        };
    },
    async mounted() {
        await this.getCurrentUser();
        await this.fetchChatUsers();
        await this.fetchFriends();

        if (this.selectedUserId) {
            const friend = this.friends.find(friend => friend.user_id === this.selectedUserId);

            if (friend) {
                this.joinRoom(friend);
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

        async getFriends() {
        // Add your logic to fetch friends here
        console.log("Fetching friends...");
        // Example: fetch friends data from API
        try {
            const response = await friendService.getFriends(this.currentUser.user_id);
            this.users = response.data;
        } catch (error) {
            console.error("Error fetching friends:", error);
        }
        },
        async fetchFriends() {
            try {
                const response = await axios.get("http://localhost:8080/friends", {
                    withCredentials: true,
                });
                this.friends = response.data;

                console.log("Fetched friends:", this.friends);
            } catch (error) {
                console.error("Error fetching friends:", error);
                this.errorMessage = "Failed to load friends. Please try again later.";
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
        async fetchChatUsers() {
            try {
                const response = await axios.get(`https://music-connect-chat-555448022527.us-central1.run.app/users/${this.currentUser.user_id}/chat-history`, {
                    withCredentials: true,
                });
                const userIds = response.data;
                console.log("Fetched user IDs:", this.selectedUserId);

                const userDetailsPromises = userIds.filter(userID => userID).map(async (userID) => {
                    try {
                        const userResponse = await UserService.getUserByFirebaseUID(userID);
                        return userResponse;
                    } catch (error) {
                        console.error(`Error fetching details for user ID ${userID}:`, error);
                        return null;
                    }
                });

                const users = await Promise.all(userDetailsPromises);
                this.users = users.filter(user => user !== null);

                console.log("Fetched chat users:", this.users);
            } catch (error) {
                console.error("Error fetching chat history:", error);
                this.errorMessage = "Failed to load chat users. Please try again later.";
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
                otherUserProfilePic: user.profile_photo_url || "/profile.svg",
                showScrollArrow: false,
            };

            this.loading = true;

            try {
                const otherUserID = roomName.split("-").filter(id => id !== this.currentUser.user_id)[0];
                const response = await this.getOtherUsers(otherUserID);
                newRoom.otherUserName = response;

                const messagesResponse = await axios.get(`https://music-connect-chat-555448022527.us-central1.run.app/rooms/${roomName}/messages`, { withCredentials: true });
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

            const ws = new WebSocket(`wss://music-connect-chat-555448022527.us-central1.run.app/ws?userid=${this.currentUser.user_id}&room=${roomName}`);

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
        }
    },
};
</script>

<style scoped>
:root {
    --primary-bg-light: #ffffff;
    --primary-bg-dark: #1f1f1f;
    --secondary-bg-light: #f5f5f5;
    --secondary-bg-dark: #333333;
    --text-light: #343a40;
    --text-dark: #f5f5f5;
    --border-light: #eee;
    --border-dark: #555555;
    --button-bg-light: #007bff;
    --button-bg-dark: #0066cc;
    --message-sent-light: #ffffff;
    --message-sent-dark: rgba(0, 123, 255, 0.5);
    --message-received-light: #ffffff;
    --message-received-dark: #444444;
}

.chat-header {
    display: flex;
    align-items: center;
    padding: 0.5rem 0;
    background-color: var(--secondary-bg-light);
    border-radius: 8px 8px 0 0;
    font-size: 1.5rem;
    font-weight: 600;
}

.user-info {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.user-avatar {
    width: 30px;   /* Smaller size */
    height: 30px;  /* Smaller size */
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid #ddd;
}

.user-name {
    font-size: 1.2rem;
    font-weight: 500;
    color: var(--text-light);
}


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
    background-color: var(--message-sent-light);
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
    margin-top: 50px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.chat-header {
    text-align: center;
    font-size: 1.5rem;
    font-weight: 600;
    padding: 0.5rem 0;
    background-color: var(--secondary-bg-light);
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
    background-color: var(--message-sent-light);
    color: var(--text-light);
    border: 1px solid var(--border-light);
}

.message-wrapper.received .message-bubble {
    background-color: var(--message-received-light);
    color: var(--text-light);
    border: 1px solid var(--border-light);
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
    background-color: var(--secondary-bg-light);
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

/* Dark Mode */
@media (prefers-color-scheme: dark) {
    body {
        background-color: var(--primary-bg-dark);
        color: var(--text-dark);
    }

    .chat-container {
        background-color: var(--primary-bg-dark);
    }

    .user-list {
        background-color: var(--primary-bg-dark);
        border: 1px solid var(--border-dark);
    }

    .user-list li:hover {
        background-color: rgba(255, 255, 255, 0.1);
    }

    
    .user-list li.active {
        background-color: rgba(255, 255, 255, 0.1);

    }

    .chat-card {
        background-color: var(--secondary-bg-dark);
        border: 1px solid var(--border-dark);
    }

    .chat-header {
        background-color: var(--secondary-bg-dark);
    }

    .chat-body {
        background-color: var(--secondary-bg-dark);
    }

    .message-wrapper.sent .message-bubble {
        background-color: var(--message-sent-dark);
        color: var(--text-dark);
        border: 1px;
        border-color: var(--border-dark);
    }

    .message-wrapper.received .message-bubble {
        background-color: var(--message-received-dark);
        color: var(--text-dark);
        border: 1px solid var(--border-dark);
    }

    .user-list li.active {
        background-color: rgba(0, 0, 0, 0.1);
    }

    .chat-footer {
        background-color: var(--secondary-bg-dark);
        color: var(--text-dark);
    }
}
</style>
