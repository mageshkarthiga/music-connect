<template>
    <Card class="user-card">
        <template #content>
            <div class="user-card-content">
    <!-- Left: Profile Picture & Username -->
    <div class="flex items-center space-x-6 p-2 min-w-0">
        <Avatar :image="user.profile_photo_url || defaultProfilePic" shape="circle" size="large" />
        <div class="truncate max-w-xs">
            <h3 class="user-name clickable-link truncate" @click="redirectToProfile">{{ user.user_name }}</h3>
            <p class="truncate text-sm text-gray-600">{{ user.email_address }}</p>
        </div>
    </div>

    <!-- Right: Action Buttons -->
    <div class="flex items-center gap-2 ml-auto">
        <Button v-if="accept" label="Accept" icon="pi pi-user-plus" @click="onAccept" severity="success" />
        <Button v-if="reject" label="Reject" icon="pi pi-times" @click="onReject" severity="danger" />
    </div>

<!-- Centered Remove button below -->
<div v-if="remove" class="w-full flex justify-center mt-2">
    <Button label="Remove" icon="pi pi-user-minus" @click="onRemove" severity="danger" />
</div>

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
    padding: 0; /* Let content manage padding */
    overflow: hidden;
    border-radius: 0.75rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.user-card-content {
    padding: 1rem 1.5rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    flex-wrap: wrap;
    width: 100%;
    box-sizing: border-box;
}

.user-card-content > div {
    min-width: 0;
}

.user-name {
    margin: 0;
    font-size: 1.125rem;
    font-weight: 600;
}

.clickable-link {
    color: black;
    cursor: pointer;
    transition: color 0.2s ease-in-out;
    display: inline-block;
}

.clickable-link:hover {
    color: #10b981;
}

.actions {
    display: flex;
    gap: 0.5rem;
    flex-wrap: wrap;
    justify-content: flex-end;
}

.p-avatar {
    flex-shrink: 0;
}

</style>