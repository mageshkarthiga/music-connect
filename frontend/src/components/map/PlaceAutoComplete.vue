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
import { ref } from "vue";
import AutoComplete from "primevue/autocomplete";
import { API_BASE_URL } from "@/service/apiConfig";
export default {
  name: "PlaceAutoComplete",
  components: { AutoComplete },
  emits: ["place-selected"],
  setup(_, { emit }) {
    const query = ref("");
    const suggestions = ref([]);

    const fetchSuggestions = async (event) => {
      const input = event.query;
      if (!input) {
        suggestions.value = [];
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
        suggestions.value = data;
      } catch (error) {
        console.error(error);
      }
    };

    const selectSuggestion = (event) => {
      const selected = event.value;
      query.value = selected;
      suggestions.value = [];
      emit("place-selected", selected);
    };

    return {
      query,
      suggestions,
      fetchSuggestions,
      selectSuggestion,
    };
  },
};
</script>
