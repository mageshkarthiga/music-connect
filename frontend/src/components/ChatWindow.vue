<template>
    <div class="chat-container">
        <Card class="chat-card">
            <template #header>
                <div class="chat-header">
                    <h3>💬 Chat</h3>
                </div>
            </template>

            <template #content>
                <div class="chat-body" ref="chatBody">
                    <div v-for="(message, index) in messages" :key="index" class="message-wrapper"
                        :class="{ 'sent': message.isSent, 'received': !message.isSent }">
                        <div class="message-bubble">
                            {{ message.message }}
                        </div>
                    </div>
                </div>
            </template>

            <template #footer>
                <div class="chat-footer">
                    <Textarea v-model="newMessage" rows="2" autoResize placeholder="Type your message..."
                        @keyup.enter.exact="sendMessage" class="chat-input" />
                    <Button icon="pi pi-send" class="p-button-rounded p-button-primary" @click="sendMessage"
                        aria-label="Send" />
                </div>
            </template>
        </Card>
    </div>
</template>

<script>
export default {
    name: "ChatWindow",
    data() {
        return {
            messages: [],
            newMessage: "",
            ws: null,
        };
    },
    mounted() {
        this.connectToWebsocket();
    },
    methods: {
        connectToWebsocket() {
            this.ws = new WebSocket("ws://localhost:1323/ws");

            this.ws.onopen = () => {
                console.log("WebSocket connection established");
            };

            this.ws.onclose = () => {
                console.log("WebSocket connection closed");
            };

            this.ws.addEventListener("message", (event) => {
                this.handleNewMessage(event);
            });
        },

        handleNewMessage(event) {
            let data = event.data;
            console.log("Received data:", data);
            data = data.split(/\r?\n/);
            for (let i = 0; i < data.length; i++) {
                if (data[i].trim() !== "") {
                    let msg = JSON.parse(data[i]);
                    this.messages.push({
                        message: msg.message,
                        isSent: false,
                    });
                }
            }
            console.log("Messages:", this.messages);
            this.scrollToBottom();
        },

        sendMessage() {
            if (this.newMessage.trim() !== "") {
                this.ws.send(JSON.stringify({ message: this.newMessage }));
                this.messages.push({
                    message: this.newMessage,
                    isSent: true,
                });
                this.newMessage = "";
                this.scrollToBottom();
            }
        },

        scrollToBottom() {
            this.$nextTick(() => {
                const chatBody = this.$refs.chatBody;
                if (chatBody) {
                    chatBody.scrollTop = chatBody.scrollHeight;
                }
            });
        }

    },
};
</script>

<style scoped>
.chat-container {
    max-width: 600px;
    margin: 0 auto;
    padding: 1rem;
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
