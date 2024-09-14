<template>
  <div class="container mx-auto p-8">
    <div v-if="game" class="flex flex-col md:flex-row gap-8">
      <div class="md:w-1/3 justify-center">
        <NuxtImg
          :src="game.cover"
          :alt="game.name + ' cover'"
          width="300"
          height="400"
          class="rounded-lg shadow-lg"
        />
      </div>
      <div class="md:w-2/3">
        <h1 class="text-3xl font-bold text-white mb-2">{{ game.name }}</h1>
        <div v-if="game.platforms" class="text-white mb-4">
          <h2 class="text-xl font-semibold mb-2">Platforms:</h2>
          <ul class="flex flex-wrap gap-2">
            <li
              v-for="platform in game.platforms"
              :key="platform.id"
              class="bg-gray-700 px-2 py-1 rounded-md text-sm"
            >
              {{ platform.name }}
            </li>
          </ul>
        </div>
        <p v-if="game.summary" class="text-white text-lg mb-8">
          {{ game.summary }}
        </p>
        <!-- Placeholder for future table -->
        <div class="bg-gray-800 p-4 rounded-lg">
          <p class="text-white text-center">
            Additional game information will be displayed here
          </p>
        </div>
      </div>
    </div>
    <p v-else class="text-white text-xl text-center">Loading...</p>
  </div>
</template>

<script lang="ts" setup>
import type { Game } from "~/server/api/apiTypes";

const gameID = useRoute().params.id as string;
const game = ref<Game | null>(null);

onMounted(async () => {
  const { data, error } = await useFetch<Game>("/api/games", {
    query: {
      gameID: gameID,
    },
  });

  if (error.value) {
    console.error("Error fetching game data:", error.value);
    return;
  }

  if (data.value) {
    game.value = data.value;
  }
});
</script>

<style></style>
