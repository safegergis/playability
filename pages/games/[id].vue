<template>
  <div class="flex flex-col items-center">
    <div class="flex items-center justify-center">
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
            <div class="flex items-center">
              <h1 class="text-3xl font-bold mb-2">
                {{ game.name }}
              </h1>
              <AccesibilityScoreBadge
                v-if="game.accessibility_score"
                :score="game.accessibility_score"
                class="ml-auto mt-4"
              />
            </div>
            <div v-if="game.platforms" class="mb-4">
              <h2 class="text-xl font-semibold mb-2">Platforms:</h2>
              <ul class="flex flex-wrap gap-2">
                <li
                  v-for="platform in platforms"
                  :key="platform.id"
                  class="bg-zinc-700 px-2 py-1 rounded-md"
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
            <p v-if="game.summary" class="text-lg mb-8">
              {{ game.summary }}
            </p>

            <Card v-if="game" class="dark p-4 rounded-lg:">
              <h2 class="text-xl font-semibold mb-4">
                Essential Accessibility Features
              </h2>
              <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="bg-zinc-700 p-4 rounded-lg">
                  <h3 class="text-lg font-medium mb-2">Colorblind Mode</h3>
                  <p class="text-gray-300">
                    {{ parseAccessibility(game.color_blind || "unknown") }}
                  </p>
                </div>
                <div class="bg-zinc-700 p-4 rounded-lg">
                  <h3 class="text-lg font-medium mb-2">Closed Captions</h3>
                  <p class="text-gray-300">
                    {{ parseAccessibility(game.closed_captions || "unknown") }}
                  </p>
                </div>
                <div class="bg-zinc-700 p-4 rounded-lg">
                  <h3 class="text-lg font-medium mb-2">Controller Support</h3>
                  <p class="text-gray-300">
                    {{
                      parseAccessibility(
                        game.full_controller_support || "unknown"
                      )
                    }}
                  </p>
                </div>
                <div class="bg-zinc-700 p-4 rounded-lg">
                  <h3 class="text-lg font-medium mb-2">
                    Full Controller Remapping
                  </h3>
                  <p class="text-gray-300">
                    {{
                      parseAccessibility(game.controller_remapping || "unknown")
                    }}
                  </p>
                </div>
              </div>
            </Card>
          </div>
        </div>
        <p v-else class="text-xl text-center">Loading...</p>
      </div>
    </div>
    <div class="container mx-auto p-8">
      <h2 class="text-3xl font-semibold mb-4">Accessibility Reports</h2>
      <hr class="w-full mb-4" />
      <ReportButton class="mb-4" :game="gameID" />
      <ReportCardGrid :game="gameID" />
    </div>
  </div>
</template>

<script lang="ts" setup>
const platformDefinitions = [
  { id: 48, name: "Playstation 4", icon: "mdi:sony-playstation" },
  { id: 49, name: "Xbox One", icon: "mdi:microsoft-xbox" },
  { id: 6, name: "PC", icon: "mdi:steam" },
  { id: 130, name: "Nintendo Switch", icon: "mdi:nintendo-switch" },
  { id: 167, name: "Playstation 5", icon: "mdi:sony-playstation" },
  { id: 168, name: "Xbox Series X", icon: "mdi:microsoft-xbox" },
];

const gameID = Number.parseInt(useRoute().params.id as string);
const game = ref<Game | null>();

const response = await useFetch<Game>("http://localhost:8080/games", {
  method: "GET",
  query: {
    id: gameID,
  },
});
console.log("response: ", response.data.value);
if (response.data.value) {
  game.value = response.data.value;
  if (
    game.value.platforms?.includes(48) &&
    game.value.platforms?.includes(167)
  ) {
    game.value.platforms = game.value.platforms.filter(
      (platform) => platform !== 48
    );
  }
  if (
    game.value.platforms?.includes(49) &&
    game.value.platforms?.includes(168)
  ) {
    game.value.platforms = game.value.platforms.filter(
      (platform) => platform !== 49
    );
  }
}
const parseAccessibility = (resp: string) => {
  if (resp === "unknown") {
    return "Data not available";
  } else if (resp === "limited") {
    return "Limited implementation";
  } else if (resp === "true") {
    return "Supported";
  } else if (resp === "false") {
    return "Not supported";
  } else {
    return "Data not available";
  }
};

const platforms = computed(() => {
  return platformDefinitions.filter((platform) =>
    game.value?.platforms?.includes(platform.id)
  );
});
</script>

<style></style>
