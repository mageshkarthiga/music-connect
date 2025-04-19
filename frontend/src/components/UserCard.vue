<template>
    <Card class="user-card">
        <div class="user-card-content">
            <img :src="user.profilePhotoUrl || defaultProfilePic" alt="Profile Picture" class="profile-pic" />
            <div class="user-info">
                <h3 class="user-name">{{ user.name }}</h3>
            </div>
        </div>
        <div class="actions">
            <Button label="Accept" icon="pi pi-check" class="p-button-success" @click="onAccept" />
            <Button label="Reject" icon="pi pi-times" class="p-button-danger" @click="onReject" />
        </div>
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
</style>