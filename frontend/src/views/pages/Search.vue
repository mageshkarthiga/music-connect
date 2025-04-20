<script>
import UserService from '@/service/UserService';
import EventService from '@/service/EventService';
import FriendService from '@/service/FriendService';
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

  // Fetch music similarity for each user
  const usersWithSimilarity = await Promise.all(usersWithStatus.map(async user => {
      const similarity = await getMusicSimilarity(user.user_id);
      return { ...user, similarity };
    }));

    allUsers.value = usersWithSimilarity;
    filteredUsers.value = [...usersWithSimilarity];
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
=======
export default {
  name: 'Search',
  data() {
    return {
      allUsers: [],
      filteredUsers: [],
      allEvents: [],
      filters1: { global: { value: null, matchMode: 'contains' } },
      isLoading: true,
      filterValue: 'desc',
      filterOptions: [
        { label: 'Highest to Lowest Compatibility', value: 'desc' },
        { label: 'Lowest to Highest Compatibility', value: 'asc' },
      ],
      currentUserId: null,
      toast: null,
      API_BASE_URL: process.env.VUE_APP_API_BASE_URL,
    };
  },
  methods: {
    async fetchData() {
>>>>>>> e68f37b405d8f0c98871a6375d0d895f6bab2897
      try {
        const meResponse = await axios.get(`${this.API_BASE_URL}/me`, { withCredentials: true });
        this.currentUserId = meResponse.data.user_id;

        const users = await UserService.getAllUsers();
        const filteredUsersList = users.filter(user => user.user_id !== this.currentUserId);

        // Fetch friendship status
        const friendsStatus = await this.fetchFriendshipStatus(filteredUsersList);

        // Map friendship status to users
        const usersWithStatus = filteredUsersList.map(user => {
          const status = friendsStatus.find(fs => fs.user_id === user.user_id);
          return { ...user, status: status ? status.status : 'none' };
        });

        // Fetch music similarity for each user
        const usersWithSimilarity = await Promise.all(usersWithStatus.map(async user => {
          const similarity = await this.getMusicSimilarity(user.user_id);
          return { ...user, similarity };
        }));

        this.allUsers = usersWithSimilarity;
        this.filteredUsers = [...usersWithSimilarity];
        this.onFilterSelect('desc');

        this.allEvents = await EventService.getAllEvents();
        this.isLoading = false;
      } catch (error) {
        console.error("Failed to fetch data:", error);
        this.isLoading = false;
      }
    },
    async fetchFriendshipStatus(users) {
      try {
        const statusPromises = users.map(async (user) => {
          try {
            const statusRes = await axios.get(
              `${this.API_BASE_URL}/friendship/${user.user_id}/status`,
              { withCredentials: true }
            );

            const userStatus = statusRes.data.status;

            // If the current user sent the request
            if (userStatus === 'pending' && statusRes.data.friend_id == this.currentUserId) {
              return {
                user_id: user.user_id,
                status: 'requested',
              };
            }

            return {
              user_id: user.user_id,
              status: userStatus,
            };
          } catch (error) {
            // Handle 404s gracefully (no friendship yet)
            if (error.response && error.response.status === 404) {
              return {
                user_id: user.user_id,
                status: 'none',
              };
            }

            // Log other errors (real issues like 500s or network errors)
            console.error(`Error fetching friendship status for user ${user.user_id}:`, error);
            return {
              user_id: user.user_id,
              status: 'none',
            };
          }
        });

        return await Promise.all(statusPromises);
      } catch (error) {
        console.error('Error fetching friendship statuses:', error);
        return [];
      }
    },
    async getMusicSimilarity(userId) {
      try {
        const response = await axios.get(
          `${this.API_BASE_URL}/calculateSimilarity?user_id1=${this.currentUserId}&user_id2=${userId}`,
          { withCredentials: true }
        );
        return response.data.similarity || 0; // Default to 0 if similarity is not present
      } catch (error) {
        console.error("Error fetching similarity:", error);
        return 0; // Default similarity score in case of error
      }
    },
    onFilterSelect(value) {
      this.filterValue = value;

      if (value === 'desc') {
        this.filteredUsers = [...this.allUsers].sort((a, b) => {
          // First, compare by similarity score (descending)
          if ((b.similarity ?? 0) !== (a.similarity ?? 0)) {
            return (b.similarity ?? 0) - (a.similarity ?? 0);  // Sort by similarity (desc)
          }

          // If similarity is the same, compare by request timestamp (for 'requested' status)
          if (a.status === 'requested' && b.status === 'requested') {
            return new Date(b.request_timestamp) - new Date(a.request_timestamp);  // Sort by most recent request
          }

          // If one user is 'requested' and the other is not, prioritize 'requested' ones
          if (a.status === 'requested') return -1;
          if (b.status === 'requested') return 1;

          return 0; // No change if both are not 'requested'
        });
      } else if (value === 'asc') {
        this.filteredUsers = [...this.allUsers].sort((a, b) => {
          // First, compare by similarity score (ascending)
          if ((a.similarity ?? 0) !== (b.similarity ?? 0)) {
            return (a.similarity ?? 0) - (b.similarity ?? 0);  // Sort by similarity (asc)
          }

          // If similarity is the same, compare by request timestamp (for 'requested' status)
          if (a.status === 'requested' && b.status === 'requested') {
            return new Date(a.request_timestamp) - new Date(b.request_timestamp);  // Sort by earliest request
          }

          // If one user is 'requested' and the other is not, prioritize 'requested' ones
          if (a.status === 'requested') return -1;
          if (b.status === 'requested') return 1;

          return 0; // No change if both are not 'requested'
        });
      },
    async addFriend(userId) {
      const user = this.allUsers.find(u => u.user_id === userId);
      if (user?.status === 'pending' || user?.status === 'accepted') {
        this.$toast.add({
          severity: 'warn',
          summary: 'Friend Request Not Sent',
          detail: `${user?.user_name ?? 'This user'} is already your friend or has a pending request.`,
        });
        return;
      }
      try {
        console.log("Sending friend request to user:", userId);
        await FriendService.sendFriendRequest(userId);
        if (user) {
          user.status = 'pending';
        }
        this.$toast.add({
          severity: 'info',
          summary: 'Friend Request Sent',
          detail: `Sent a follow request to ${user?.user_name ?? 'this user'}`,
        });
      } catch (error) {
        if (this.$toast) {
          if (error.response && error.response.status === 409) {
            if (user) {
              user.status = 'pending';
            }
            this.$toast.add({
              severity: 'warn',
              summary: 'Already Sent',
              detail: `${user?.user_name ?? 'This user'} already has a pending or accepted request.`,
            });
          } else {
            console.error("Failed to send friend request:", error);
            if (this.$toast) {
              this.$toast.add({
                severity: 'error',
                summary: 'Error',
                detail: 'Failed to send friend request. Please try again.',
              });
            }
          }
        }
      }
      
      // If one user is 'requested' and the other is not, prioritize 'requested' ones
      if (a.status === 'requested') return -1;
      if (b.status === 'requested') return 1;

      return 0; // No change if both are not 'requested'
    });
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
    await axios.post(`http://localhost:8080/friends/${userId}/request`, {}, { withCredentials: true });
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
    await axios.post(`http://localhost:8080/friends/${userId}/accept`, {}, { withCredentials: true });
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
    await axios.post(`http://localhost:8080/friends/${userId}/reject`, {}, { withCredentials: true });
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


=======
    },
    async acceptRequest(userId) {
      const user = this.allUsers.find(u => u.user_id === userId);
      try {
        await FriendService.acceptFriendRequest(userId);
        this.$toast.add({
          severity: 'info',
          summary: 'Friend Request Accepted',
          detail: `You are now friends with ${user?.user_name ?? 'this user'}`,
        });
      } catch (error) {
        console.error("Failed to accept friend request:", error);
      }
    },
    async rejectRequest(userId) {
      const user = this.allUsers.find(u => u.user_id === userId);
      try {
        await FriendService.rejectFriendRequest(userId);
        this.$toast.add({
          severity: 'info',
          summary: 'Friend Request Rejected',
          detail: `You rejected the friend request from ${user?.user_name ?? 'this user'}`,
        });
      } catch (error) {
        console.error("Failed to reject friend request:", error);
      }
    },
  },
  mounted() {
    this.fetchData();
  },
};
>>>>>>> e68f37b405d8f0c98871a6375d0d895f6bab2897
</script>

<template>

    
  <div class="card">
    <Toast ref="toast" />

    <Tabs>
      <TabPanel header="Users">



        <DataTable :value="filteredUsers" :paginator="true" :rows="10" dataKey="user_id" :filters="filters1"
          filterDisplay="menu" :loading="isLoading" :globalFilterFields="['user_name']">
          <template #header>
            <div class="w-full flex items-center justify-between mb-4 mt-8">
              <!-- Header Text aligned to the left -->
              <div class="font-semibold text-xl flex-1">Your Friends & Requests </div>

              <!-- Filter Buttons aligned to the right -->
              <div class="flex items-center gap-2 ml-auto">
                <Button label="Sort by Music Similarity" icon="pi pi-sort-alt"
                  class="p-button-outlined p-button-secondary mr-2" @click="onFilterSelect('desc')" />
                <Button label="Sort by Incoming Requests" icon="pi pi-sort" class="p-button-outlined p-button-secondary"
                  @click="onFilterSelect('asc')" />
              </div>

              <!-- Search Bar aligned to the top-right -->
              <div class="absolute top-0 right-0 flex items-center gap-2">
                <i class="pi pi-search text-gray-500" />
                <InputText v-model="filters1.global.value" placeholder="Search for user" class="w-[18rem]" />
              </div>
            </div>
          </template>

          <div class="flex justify-end mb-4">

          </div>

          <Column header="Users" style="min-width: 25rem ">
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <template v-if="isLoading">
                  <Skeleton />
                </template>
                <template v-else>
                  <Avatar :image="data.profile_photo_url || '/profile.svg'" shape="circle" size="large" />
                  <span class="clickable-link">
                    <a :href="`/profile?user_id=${data.user_id}`">
                      {{ data.user_name }}
                    </a>
                    <div class="text-xs font-bold" :class="{
                      'text-green-500': (data.similarity ?? 0) >= 0.6,  // Green for high similarity
                      'text-yellow-500': (data.similarity ?? 0) >= 0.3 && (data.similarity ?? 0) < 0.6,  // Yellow for medium similarity
                      'text-red-500': (data.similarity ?? 0) < 0.3  // Red for low similarity
                    }">
                      Match: {{ ((data.similarity ?? 0) * 100).toFixed(2) }}%
                    </div>
                  </span>
                </template>
              </div>
            </template>
          </Column>

          <Column>
            <template #body="{ data }">
              <div class="flex items-center gap-2">
                <template v-if="isLoading">
                  <Skeleton />
                </template>
                <template v-else>
                  <Button v-if="data.status === 'pending'" label="Pending"
                    class="p-button-outlined p-button-secondary p-button-sm" :disabled="true" />

                  <span v-else-if="data.status === 'accepted'">

                    <RouterLink :to="{ name: 'chat', params: { user_id: data.firebase_uid } }">

                      <Button label="Chat" icon="pi pi-comments" class="p-button-sm p-button-success"
                        @click.stop="goToChat(data.user_id)" />

                    </RouterLink>
                  </span>

                  <span class='flex gap-2' v-else-if="data.status === 'requested'">
                    <Button label="Accept" icon="pi pi-check" class="p-button-sm p-button-success"
                      @click.stop="acceptRequest(data.user_id)" />

                    <Button label="Reject" icon="pi pi-times" class="p-button-sm p-button-danger"
                      @click.stop="rejectRequest(data.user_id)" />
                  </span>

                  <Button v-else label="Add Friend" icon="pi pi-user-plus" class="p-button-outlined p-button-sm"
                    @click.stop="addFriend(data.user_id)" />
                </template>
              </div>
            </template>
          </Column>
        </DataTable>
      </TabPanel>

      <!-- Your events tab here -->
    </Tabs>
  </div>
<<<<<<< HEAD
</template>
=======
</template>
