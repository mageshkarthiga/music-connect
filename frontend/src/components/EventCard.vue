<template>
  <div class="max-w-[300px] min-w-[250px] flex-shrink-0 cursor-pointer relative" @click="openEventUrl">
    <Card class="w-full h-full">
      <template v-slot:title>
        <div class="flex items-center justify-between mb-0">
          <div class="font-semibold text-xl mb-4 line-clamp-3 min-h-[50px]">
            {{ event.event_name }}
          </div>
          <button @click.stop="toggleLike" class="hover:scale-110 transition-transform duration-200">
            <component :is="'lucide-heart'" class="w-5 h-5 transition-all duration-200"
              :class="isLiked ? 'text-red-500 fill-red-500' : 'text-gray-400 hover:text-red-400'" />
          </button>
        </div>
      </template>

      <template v-slot:content>
        <img :src="event.event_image_url" alt="Event image" class="w-full h-40 object-cover rounded-xl mb-3" />
        <p class="text-sm text-gray-800 line-clamp-2 mt-1 mb-4"> 
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
    liked: Boolean,
  },
  data() {
    return {
      isLiked: this.liked,
    };
  },
  watch: {
    liked(newVal) {
      this.isLiked = newVal;
    }
  },
  methods: {
    openEventUrl() {
      if (this.event.event_url) {
        window.open(this.event.event_url, "_blank");
      }
    },
    async toggleLike() {
      try {
        if (this.isLiked) {
          await EventService.unlikeEvent(this.event.event_id);
          this.isLiked = false;
          this.$emit("event-unliked", this.event.event_id);
        } else {
          await EventService.likeEvent(this.event.event_id);
          this.isLiked = true;
          this.$emit("event-liked", this.event.event_id);
        }
      } catch (err) {
        console.error("Failed to toggle like:", err);
      }
    },
  },
};
</script>

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
