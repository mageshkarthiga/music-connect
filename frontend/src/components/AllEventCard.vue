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
                :is="liked ? 'lucide-heart' : 'lucide-heart'"
                class="w-5 h-5"
                :class="liked ? 'text-red-500 fill-red-500' : 'text-gray-400'"
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
    data() {
      return {
        liked: false,
      };
    },
    methods: {
      openEventUrl() {
        if (this.event?.event_url) {
          window.open(this.event.event_url, "_blank");
        }
      },
      async toggleLike() {
        if (this.liked) return;
        try {
          await EventService.likeEvent(this.event.event_id);
          this.liked = true;
        } catch (err) {
          console.error("Failed to like event:", err);
        }
      },
    },
  };
  </script>
  