<template>
    <AppTopbar />
    <ChatWindow :selectedUserId="userId" />
</template>

<script>
import { onMounted } from 'vue';
import { useRoute } from 'vue-router';
import ChatWindow from '@/components/ChatWindow.vue';
import AppTopbar from '@/layout/AppTopbar.vue';

export default {
    components: {
        ChatWindow,
        AppTopbar
    },
    setup() {
        const route = useRoute();
        const userId = route.params.user_id || null;

        onMounted(() => {
            if (userId) {
                console.log(`Opening chat with user_id: ${userId}`);
            }
        });

        return {
            userId,
        };
    },
};
</script>