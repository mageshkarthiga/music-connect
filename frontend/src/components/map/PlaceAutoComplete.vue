<template>
  <AutoComplete
    v-model="query"
    :suggestions="suggestions"
    @complete="fetchSuggestions"
    @select="selectSuggestion"
    @type="selectSuggestion"
    placeholder="Enter an address"
  />
</template>

<script>
import AutoComplete from "primevue/autocomplete";

export default {
  name: "PlaceAutoComplete",
  components: { AutoComplete },
  emits: ["place-selected"],
  data() {
    return {
      query: "",
      suggestions: [],
    };
  },
  methods: {
    async fetchSuggestions(event) {
      const input = event.query;
      if (!input) {
        this.suggestions = [];
        return;
      }
      try {
        const response = await fetch(
          `${API_BASE_URL}/places/autocomplete?input=${encodeURIComponent(input)}`
        );
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        const data = await response.json();
        this.suggestions = data;
      } catch (error) {
        console.error(error);
      }
    },
    selectSuggestion(event) {
      const selected = event.value;
      this.query = selected;
      this.suggestions = [];
      this.$emit("place-selected", selected);
    },
  },
};
</script>
