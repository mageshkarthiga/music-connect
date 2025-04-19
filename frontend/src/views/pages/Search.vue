<script setup>
import { FilterMatchMode } from '@primevue/core/api';
import { onBeforeMount, ref, watch } from 'vue';
import UserService from '@/service/UserService';
import EventService from '@/service/EventService';
import axios from 'axios';

const allUsers = ref([]);
const filteredUsers = ref([]);
const allEvents = ref([]);
const filters1 = ref({ global: { value: null, matchMode: FilterMatchMode.CONTAINS } }); // Initializing properly
const isLoading = ref(true);
const filterValue = ref('desc'); // Default to 'desc' (Highest to Lowest)
const filterOptions = ref([
  { label: 'Highest to Lowest Compatibility', value: 'desc' },
  { label: 'Lowest to Highest Compatibility', value: 'asc' },
]);

const currentUserId = ref(null);

onBeforeMount(async () => {
  try {
    // Fetch current user
    const meResponse = await axios.get("http://localhost:8080/me", {
      withCredentials: true,
    });
    currentUserId.value = meResponse.data.user_id;

    // Fetch all users
    const users = await UserService.getAllUsers();
    const filteredUsersList = users.filter(user => user.user_id !== currentUserId.value);

    // Calculate similarity for each user
    const usersWithSimilarity = await Promise.all(filteredUsersList.map(async (user) => {
      try {
        const simResponse = await axios.get("http://localhost:8080/calculateSimilarity", {
          params: {
            user_id1: currentUserId.value,
            user_id2: user.user_id,
          },
        });
        return {
          ...user,
          similarity: simResponse.data.similarity ?? 0,
        };
      } catch (error) {
        console.error(`Failed to calculate similarity for user ${user.user_id}`, error);
        return { ...user, similarity: 0 };
      }
    }));

    allUsers.value = usersWithSimilarity;
    filteredUsers.value = [...usersWithSimilarity]; // Initialize filteredUsers with all users

    // Trigger sorting by default to highest compatibility (desc)
    onFilterSelect('desc'); // Sort from highest to lowest compatibility by default

    // Fetch all events
    allEvents.value = await EventService.getAllEvents();

    isLoading.value = false;
  } catch (error) {
    console.error("Failed to fetch data:", error);
    isLoading.value = false;
  }
});

// Handle sorting based on dropdown selection
function onFilterSelect(value) {
  filterValue.value = value;

  // Sort users based on similarity score
  if (value === 'desc') {
    filteredUsers.value = [...allUsers.value].sort((a, b) => b.similarity - a.similarity);
  } else if (value === 'asc') {
    filteredUsers.value = [...allUsers.value].sort((a, b) => a.similarity - b.similarity);
  }
}
</script>

<template>
  <div class="card">

    <Tabs>
      <TabPanel header="Users">
        <div class="flex justify-between items-center mb-2">

        </div>
        <DataTable
          :value="filteredUsers"
          :paginator="true"
          :rows="10"
          dataKey="user_id"
          :filters="filters1"
          filterDisplay="menu"
          :loading="isLoading"
          :globalFilterFields="['user_name']"
        >
          <template #header>

            
            <div class="w-full flex items-center justify-between mb-4">
  <div class="font-semibold text-xl">Search for Users or Events</div>

  <div class="flex items-center gap-2">
    <i class="pi pi-search text-gray-500" />
    <InputText
      v-model="filters1.global.value"
      placeholder="Search for user"
      class="w-[18rem]"
    />
  </div>
</div>





          </template>
          <template #empty>No users found.</template>
          <template #loading>Loading user data. Please wait.</template>
          <Column header="Users" style="min-width: 14rem">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar :image="data.profile_photo_url || '/public/profile.svg'" shape="circle" size="large" />
                <span>
                  <a :href="`/profile?user_id=${data.user_id}`" class="clickable-link">
                    {{ data.user_name }}
                  </a>
                  <div class="text-xs text-gray-500">
                    Match: {{ (data.similarity ?? 0).toFixed(2) }}
                  </div>
                </span>
              </div>
            </template>
          </Column>
          <Column>
            <template #body="{ data }">
              <RouterLink :to="{ name: 'chat', params: { user_id: data.firebase_uid } }">
                <Button label="Chat" icon="pi pi-send" class="p-button-outlined p-button-sm" />
              </RouterLink>
            </template>
          </Column>
        </DataTable>
      </TabPanel>

      <TabPanel header="Events">
        <DataTable
          :value="allEvents"
          :paginator="true"
          :rows="10"
          dataKey="event_id"
          :loading="isLoading"
          :globalFilterFields="['event_name']"
        >
          <template #header>
            <span class="p-input-icon-left">
              <i class="pi pi-search" />
              <InputText
                v-model="filters1.global.value"
                placeholder="Search for event"
                class="w-full"
              />
            </span>
          </template>
          <template #empty>No events found.</template>
          <template #loading>Loading event data. Please wait.</template>
          <Column header="Events" style="min-width: 14rem">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <Avatar :image="data.event_image_url" size="large" />
                <span>
                  <a :href="data.event_url" target="_blank" rel="noopener noreferrer" class="clickable-link">
                    {{ data.event_name }}
                  </a>
                </span>
              </div>
            </template>
          </Column>
        </DataTable>
      </TabPanel>
    </Tabs>
  </div>
</template>

<style scoped>
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
