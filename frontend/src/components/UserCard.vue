<template>
    <Card class="user-card">
        <template #content>
            <div class="user-card-content">
                <Avatar :image="user.profile_photo_url || defaultProfilePic" shape="circle" size="large" />
                <div class="user-info">
                    <h3 class="user-name clickable-link" v-on:click="redirectToProfile">{{ user.user_name }}</h3>
                    <p>&nbsp;{{ user.email_address }}</p>
                </div>
            </div>
            <div class="actions">
                <Button v-if="accept" label="Accept" icon="pi pi-user-plus" @click="onAccept" severity="success"/>
                <Button v-if="reject" label="Reject" icon="pi pi-times" @click="onReject" severity="danger"/>
                <Button v-if="remove" label="Remove" icon="pi pi-user-minus"  @click="onRemove" severity="danger" />
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
        accept: {
            type: Boolean,
            default: false,
        },
        reject: {
            type: Boolean,
            default: false,
        },
        remove: {
            type: Boolean,
            default: false,
        },
    },
    data() {
        return {
            defaultProfilePic: "/profile.svg", 
        };
    },
    methods: {
        onAccept() {
            this.$emit("accept", this.user);
        },
        onReject() {
            this.$emit("reject", this.user);
        },
        onRemove() {
            this.$emit("remove", this.user);
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
    justify-content: center; 
    align-items: center; 
    gap: 0.5rem; 
    width: 100%; 
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