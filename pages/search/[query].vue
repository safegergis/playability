<template>
  <div class="flex flex-col items-center justify-center min-h-screen p-4">
    <h1 class="text-3xl font-bold text-foreground mb-6">
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
          @click="onSearchResultClick(result.id)"
        >
          <Card
            class="bg-zinc-800 rounded-lg p-4 shadow-md transition ease-in-out delay-150 hover:scale-105"
          >
            <h2 class="text-xl font-semibold text-foreground">
              {{ result.name }}
            </h2>
          </Card>
        </li>
      </ul>
    </div>
    <div
      v-else-if="searchResults && searchResults.length === 0"
      class="text-white text-xl"
    >
      No results found for "{{ searchQuery }}"
    </div>
    <div v-if="status === 'pending'" class="text-foreground text-xl">
      Loading...
    </div>
  </div>
</template>

<script lang="ts" setup>
const searchQuery = useRoute().params.query as string;

const { data, status } = await useFetch<SearchResult[]>(
  "http://localhost:8080/search",
  {
    method: "GET",
    query: {
      search: searchQuery,
    },
  }
);

const searchResults = ref<SearchResult[]>(data.value ?? []);

const onSearchResultClick = (id: number) => {
  navigateTo(`/games/${id}`);
};
</script>

<style></style>
