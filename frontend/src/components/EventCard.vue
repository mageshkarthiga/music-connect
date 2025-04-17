<template>
  <div
    class="min-w-[250px] max-w-[300px] flex-shrink-0 cursor-pointer"
    @click="openEventUrl"
  >
    <Card class="w-full h-full">
      <template v-slot:title>
        <div class="flex items-center justify-between mb-0">
          <div class="font-semibold text-xl mb-4 clickable-link">{{ event.event_name }}</div>
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
export default {
  props: {
    event: Object,
  },
  methods: {
    openEventUrl() {
      if (this.event?.event_url) {
        window.open(this.event.event_url, "_blank");
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
