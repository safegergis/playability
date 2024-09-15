<template>
  <div class="min-h-screen flex items-center justify-center">
    <div class="container mx-auto p-8">
      <div v-if="game" class="flex flex-col md:flex-row gap-8">
        <div class="md:w-1/3 justify-center">
          <NuxtImg
            :src="game.cover_art"
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
                v-for="platform in platforms"
                :key="platform.id"
                class="bg-gray-700 px-2 py-1 rounded-md"
              >
                <Icon
                  :name="
                    game.platforms.includes(platform.id) ? platform.icon : ''
                  "
                  size="24"
                />
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
  </div>
</template>

<script lang="ts" setup>
import type { Game } from "~/server/api/apiTypes";

const platformDefinitions = [
  { id: 48, name: "Playstation", icon: "mdi:sony-playstation" },
  { id: 49, name: "Xbox", icon: "mdi:microsoft-xbox" },
  { id: 6, name: "PC", icon: "mdi:steam" },
  { id: 130, name: "Nintendo Switch", icon: "mdi:nintendo-switch" },
];

const gameID = useRoute().params.id as string;
const game = ref<Game | null>();

const response = await useFetch<Game[]>(
  `http://localhost:8080/games/${gameID}`,
  {
    method: "GET",
  }
);
console.log("response: ", response.data.value);
if (response) {
  game.value = response.data.value[0];
}

const platforms = computed(() => {
  return platformDefinitions.filter((platform) =>
    game.value?.platforms?.includes(platform.id)
  );
});

/*
if (response.error) {
  error.value = {
    statusCode: response.error.statusCode,
    statusMessage: response.error.statusMessage,
    message: `Failed to load game data: ${response.error.message}`,
  };
} else if (!response) {
  error.value = {
    statusCode: 404,
    statusMessage: "Not Found",
    message: "Game not found",
  };
  */
</script>

<style></style>
