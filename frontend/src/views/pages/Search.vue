<script setup>
import { ref, onBeforeMount } from 'vue';
import { useRouter } from 'vue-router';
import UserService from '@/service/UserService';
import EventService from '@/service/EventService';
import friendService from '@/service/FriendService';
import axios from 'axios';
import API_BASE_URL from '@/service/apiConfig';

const allUsers = ref([]);
const filteredUsers = ref([]);
const allEvents = ref([]);
const filters1 = ref({ global: { value: null, matchMode: 'contains' } });

const isLoading = ref(true);
const filterValue = ref('desc');
const filterOptions = ref([
  { label: 'Highest to Lowest Compatibility', value: 'desc' },
  { label: 'Lowest to Highest Compatibility', value: 'asc' },
]);
const currentUserId = ref(null);
const toast = ref(null);
const router = useRouter();

onBeforeMount(async () => {
  try {
    const meResponse = await axios.get("http://localhost:8080/me", { withCredentials: true });
    currentUserId.value = meResponse.data.user_id;

    const users = await UserService.getAllUsers();
    const filteredUsersList = users.filter(user => user.user_id !== currentUserId.value);

    // Fetch friendship status
    const friendsStatus = await fetchFriendshipStatus(filteredUsersList);

    // Map friendship status to users
    const usersWithStatus = filteredUsersList.map(user => {
      const status = friendsStatus.find(fs => fs.user_id === user.user_id);
      return { ...user, status: status ? status.status : 'none' };
    });

    allUsers.value = usersWithStatus;
    filteredUsers.value = [...usersWithStatus];
    onFilterSelect('desc');

    allEvents.value = await EventService.getAllEvents();
    isLoading.value = false;
  } catch (error) {
    console.error("Failed to fetch data:", error);
    isLoading.value = false;
  }
});

async function fetchFriendshipStatus(users) {
  try {
    const statusPromises = users.map(async (user) => {
      try {
        const statusRes = await axios.get(`http://localhost:8080/friendship/${user.user_id}/status`, { withCredentials: true });
        if (statusRes.status === 200) {
          const userStatus = statusRes.data.status;

          console.log(`User ID: ${user.user_id}, Status: ${userStatus}`);

          console.log (statusRes.data);
          
          // Check if the status is 'pending' and the request is sent by the current user
          if (userStatus === 'pending' && statusRes.data.friend_id == currentUserId.value) {
            return {
              user_id: user.user_id,
              status: 'requested', // Show 'requested' when the current user has sent a request
            };
          }

          return {
            user_id: user.user_id,
            status: userStatus, // Return other statuses as is
          };
        }
        return { user_id: user.user_id, status: 'none' };
      } catch (error) {
        console.error('Error fetching friendship status:', error);
        return { user_id: user.user_id, status: 'none' };
      }
    });

    return await Promise.all(statusPromises);
  } catch (error) {
    console.error('Error fetching friendship statuses:', error);
    return [];
  }
}



function onFilterSelect(value) {
  filterValue.value = value;

  if (value === 'desc') {
    filteredUsers.value = [...allUsers.value].sort((a, b) => b.similarity - a.similarity);
  } else if (value === 'asc') {
    filteredUsers.value = [...allUsers.value].sort((a, b) => a.similarity - b.similarity);
  }
}

async function addFriend(userId) {
  const user = allUsers.value.find(u => u.user_id === userId);
  if (user?.status === 'pending' || user?.status === 'accepted') {
    toast.value.add({
      severity: 'warn',
      summary: 'Friend Request Not Sent',
      detail: `${user?.user_name ?? 'This user'} is already your friend or has a pending request.`,
    });
    return;
  }

  try {
    await axios.post(`http://localhost:8080/friend/${userId}/request`, {}, { withCredentials: true });
    if (user) {
      user.status = 'pending';
    }

    toast.value.add({
      severity: 'info',
      summary: 'Friend Request Sent',
      detail: `Sent a follow request to ${user?.user_name ?? 'this user'}`,
    });
  } catch (error) {
    if (error.response && error.response.status === 409) {
      if (user) {
        user.status = 'pending';
      }
      toast.value.add({
        severity: 'warn',
        summary: 'Already Sent',
        detail: `${user?.user_name ?? 'This user'} already has a pending or accepted request.`,
      });
    } else {
      console.error("Failed to send friend request:", error);
    }
  }
}

async function acceptRequest(userId) {
  const user = allUsers.value.find(u => u.user_id === userId);
  try {
    await axios.post(`http://localhost:8080/friend/${userId}/accept`, {}, { withCredentials: true });
    if (user) {
      user.status = 'accepted';
    }

    toast.value.add({
      severity: 'info',
      summary: 'Friend Request Accepted',
      detail: `You are now friends with ${user?.user_name ?? 'this user'}`,
    });
  } catch (error) {
    console.error("Failed to accept friend request:", error);
  }
}

async function rejectRequest(userId) {
  const user = allUsers.value.find(u => u.user_id === userId);
  try {
    await axios.post(`http://localhost:8080/friend/${userId}/reject`, {}, { withCredentials: true });
    if (user) {
      user.status = 'none'; // Reset status to "none"
    }

    toast.value.add({
      severity: 'info',
      summary: 'Friend Request Rejected',
      detail: `You rejected the friend request from ${user?.user_name ?? 'this user'}`,
    });
  } catch (error) {
    console.error("Failed to reject friend request:", error);
  }
}


</script>

<template>
  <div class="card">
    <Toast ref="toast" />

    <Tabs>
      <TabPanel header="Users">
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

          <template #loading>
            <div class="rounded border p-6">
              <ul class="m-0 p-0 list-none">
                <li class="mb-4" v-for="i in 4" :key="i">
                  <div class="flex">
                    <Skeleton shape="circle" size="4rem" class="mr-2"></Skeleton>
                    <div class="self-center" style="flex: 1">
                      <Skeleton width="100%" class="mb-2"></Skeleton>
                      <Skeleton width="75%"></Skeleton>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
          </template>

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
              <div class="flex items-center gap-2">
                <Button
                  v-if="data.status === 'pending'"
                  label="Pending"
                  class="p-button-outlined p-button-secondary"
                  :disabled="true"
                />

                <span v-else-if="data.status === 'accepted'">

                <RouterLink :to="{ name: 'chat', params: { user_id: data.firebase_uid } }">

                    <Button
                    label="Chat"
                    icon="pi pi-comments"
                    class="p-button-sm p-button-success"
                    @click.stop="goToChat(data.user_id)"
                    />

                </RouterLink>
                </span>

                <span class='flex gap-2' v-else-if="data.status === 'requested' ">
                <Button
                label="Accept"
                icon="pi pi-check"
                class="p-button-sm p-button-success"
                @click.stop="acceptRequest(data.user_id)"
                />

                <Button
                label="Reject"
                icon="pi pi-times"
                class="p-button-sm p-button-danger"
                @click.stop="rejectRequest(data.user_id)"
                />
                </span>

                <Button
                  v-else
                  label="Add Friend"
                  icon="pi pi-user-plus"
                  class="p-button-outlined p-button-sm"
                  @click.stop="addFriend(data.user_id)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </TabPanel>
      
      <!-- Your events tab here -->
    </Tabs>
  </div>
</template>
