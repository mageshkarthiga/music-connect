<script setup>
import TabView from 'primevue/tabview'
import { FilterMatchMode } from '@primevue/core/api';
import { onBeforeMount, ref } from 'vue';
import UserService from '@/service/UserService';
import EventService from '@/service/EventService';

const allUsers = ref(null);
const allEvents = ref(null);
const filters1 = ref(null);
const isLoading = ref(null);

onBeforeMount(() => {
    UserService.getAllUsers().then((data) => {
        console.log('data', data);
        allUsers.value = data;
        isLoading.value = false;
    });

    EventService.getAllEvents().then((data) => {
        allEvents.value = data;
        isLoading.value = false;
    });

    initFilters1();
});

function initFilters1() {
    filters1.value = {
        global: { value: null, matchMode: FilterMatchMode.STARTS_WITH },
    };
}

</script>

<template>
    <div class="card">
        <div class="font-semibold text-xl mb-4">Search for Users or Events</div>
        <TabView>
            <TabPanel header="Users">
                <DataTable :value="allUsers" :paginator="true" :rows="10" dataKey="user_id" :rowHover="true"
                    v-model:filters="filters1" filterDisplay="menu" :loading="isLoading"
                    :globalFilterFields="['user_name']">
                    <template #header>
                        <IconField>
                            <InputIcon>
                                <i class="pi pi-search" />
                            </InputIcon>
                            <InputText v-model="filters1['global'].value" placeholder="Search for user"
                                class="w-full" />
                        </IconField>
                    </template>
                    <template #empty> No users found. </template>
                    <template #loading> Loading user data. Please wait. </template>
                    <Column header="Users" style="min-width: 14rem">
                        <template #body="{ data }">
                            <div class="flex items-center gap-2">
                                <Avatar :image="data.profile_photo_url || '/public/profile.svg'" shape="circle"
                                    size="large" />
                                <span>
                                    <a :href="`/profile?user_id=${data.user_id}`" class="clickable-link">
                                        {{ data.user_name }}
                                    </a>
                                </span>
                            </div>
                        </template>
                    </Column>
                    <Column>
                        <template #body="{ data }">
                            <RouterLink :to="{ name: 'chat', params: { user_id: data.firebase_uid } }">
                                <Button label="Chat" icon="pi pi-send" class="p-button-outlined p-button-sm"/>
                            </RouterLink>
                        </template>
                    </Column>
                </DataTable>
            </TabPanel>
            <TabPanel header="Events">
                <DataTable :value="allEvents" :paginator="true" :rows="10" dataKey="event_id" :rowHover="true"
                    :loading="isLoading" :globalFilterFields="['event_name']">
                    <template #header>
                        <IconField>
                            <InputIcon>
                                <i class="pi pi-search" />
                            </InputIcon>
                            <InputText v-model="filters1['global'].value" placeholder="Search for event"
                                class="w-full" />
                        </IconField>
                    </template>
                    <template #empty> No events found. </template>
                    <template #loading> Loading event data. Please wait. </template>
                    <Column header="Events" style="min-width: 14rem">
                        <template #body="{ data }">
                            <div class="flex items-center gap-2">
                                <!-- <img :alt="data.event_image_url" :src="data.event_image_url" style="width: 32px" /> -->
                                <Avatar :image="data.event_image_url" size="large" />
                                <span>
                                    <a :href="data.event_url" target="_blank" rel="noopener noreferrer"
                                        class="clickable-link">
                                        {{ data.event_name }}
                                    </a>
                                </span>
                            </div>
                        </template>
                    </Column>
                </DataTable>
            </TabPanel>
        </TabView>
    </div>
</template>

<style scoped lang="scss">
:deep(.p-datatable-frozen-tbody) {
    font-weight: bold;
}

:deep(.p-datatable-scrollable .p-frozen-column) {
    font-weight: bold;
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
