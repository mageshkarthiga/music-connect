<template>
  <div
    class="min-w-[250px] max-w-[300px] flex-shrink-0 cursor-pointer relative"
    @click="openEventUrl"
  >
    <Card class="w-full h-full">
      <template v-slot:title>
        <div class="flex items-center justify-between mb-0">
          <div class="font-semibold text-xl mb-4">{{ event.event_name }}</div>
          <button @click.stop="toggleLike">
            <component
              :is="'lucide-heart'"
              class="w-5 h-5"
              :class="isLiked ? 'text-red-500 fill-red-500' : 'text-gray-400'"
            />

          </button>
        </div>
      </template>

      <template v-slot:content>
        <img
          :src="event.event_image_url"
          alt="Event image"
          class="w-full h-40 object-cover rounded-xl mb-3"
        />
        <p class="text-sm text-gray-800 line-clamp-3 mt-1 mb-4">
          {{ event.event_description }}
        </p>
      </template>
    </Card>
  </div>
</template>

<script>
import { Heart } from "lucide-vue-next";
import EventService from "@/service/EventService";

export default {
  components: {
    'lucide-heart': Heart,
  },
  props: {
    event: Object,
  },
  props: {
  event: Object,
  liked: Boolean,
},
data() {
  return {
    isLiked: this.liked,
  };
},
methods: {
  async toggleLike() {
    if (this.isLiked) return; // Already liked, don’t try again
    try {
      await EventService.likeEvent(this.event.event_id);
      this.isLiked = true;
      this.$emit("event-liked", this.event.event_id); // Emit to parent
    } catch (err) {
      if (err.response?.data === "Event already liked") {
        this.isLiked = true; // Set liked locally
        return; // No toast needed
      }
      console.error("Failed to like event:", err);
    }
  },
},

    async toggleLike() {
  if (this.liked) return;
  try {
    await EventService.likeEvent(this.event.event_id);
    this.liked = true;

    // Emit to parent
    this.$emit("event-liked", this.event.event_id);
  } catch (err) {
    console.error("Failed to like event:", err);
  }
}
};
</script>
