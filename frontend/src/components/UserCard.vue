<template>
    <Card class="user-card">
        <template #content>
            <div class="user-card-content">
                <Avatar :image="user.profile_photo_url || defaultProfilePic" shape="circle" size="large" />
                <div class="user-info">
                    <h3 class="user-name clickable-link" v-on:click="redirectToProfile">{{ user.user_name }}</h3>
                </div>
            </div>
            <div class="actions">
                <Button label="Accept" icon="pi pi-check" class="p-button-success" @click="onAccept" />
                <Button label="Reject" icon="pi pi-times" class="p-button-danger" @click="onReject" />
            </div>
        </template>
    </Card>
</template>

<script>
export default {
    name: "UserCard",
    props: {
        user: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            defaultProfilePic: "/profile.svg", // Fallback profile picture
        };
    },
    methods: {
        onAccept() {
            this.$emit("accept", this.user);
        },
        onReject() {
            this.$emit("reject", this.user);
        },
        redirectToProfile() {
            const profileUrl = this.$router.resolve({ name: "profile", query: { user_id: this.user.user_id } }).href;
            window.open(profileUrl, "_blank");
        }
    },
};
</script>

<style scoped>
.user-card {
    margin-bottom: 1rem;
    padding: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
}

.user-card-content {
    display: flex;
    align-items: center;
    margin-bottom: 1rem;
}

.profile-pic {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    object-fit: cover;
    margin-right: 1rem;
}

.user-info {
    flex-grow: 1;
}

.user-name {
    margin: 0;
    font-size: 1.2rem;
    font-weight: bold;
}

.actions {
    display: flex;
    gap: 0.5rem;
}

.clickable-link {
    color: black;
    cursor: pointer;
    transition: color 0.2s ease-in-out;
}

.clickable-link:hover {
    color: #10b981;
    text-decoration: none;
}
</style>