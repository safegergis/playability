<template>
  <div class="flex flex-col items-center justify-center min-h-screen p-4">
    <h1 class="text-3xl font-bold text-white mb-6">
      Search Results for "{{ searchQuery }}"
    </h1>
    <div
      v-if="searchResults && searchResults.length > 0"
      class="w-full max-w-2xl"
    >
      <ul class="space-y-4">
        <li
          v-for="result in searchResults"
          :key="result.id"
          class="bg-stone-800 rounded-lg p-4 shadow-md"
          @click="onSearchResultClick(result.id)"
        >
          <h2 class="text-xl font-semibold text-white">{{ result.name }}</h2>
        </li>
      </ul>
    </div>
    <div
      v-else-if="searchResults && searchResults.length === 0"
      class="text-white text-xl"
    >
      No results found for "{{ searchQuery }}"
    </div>
    <div v-if="status === 'pending'" class="text-white text-xl">Loading...</div>
  </div>
</template>

<script lang="ts" setup>
import type { SearchResult } from "~/server/api/apiTypes";
const searchQuery = useRoute().params.query as string;

const { data, status } = await useFetch<string>(
  "http://localhost:8080/search",
  {
    method: "GET",
    query: {
      search: searchQuery,
    },
  }
);
const searchResults = ref<SearchResult[]>(JSON.parse(data.value as string));
const onSearchResultClick = (id: number) => {
  navigateTo(`/games/${id}`);
};
</script>

<style></style>
