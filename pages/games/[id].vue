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
              <h1 class="text-5xl font-bold mb-2">
                {{ game.name }}
              </h1>
              <AccesibilityScoreBadge :score="score" class="ml-auto mt-4" />
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
            <FeatureCard :feature-stats="featureStats" :game="game" />
          </div>
        </div>
        <p v-else class="text-xl text-center">Loading...</p>
      </div>
    </div>
    <div class="container mx-auto p-8">
      <h2 class="text-3xl font-semibold mb-4">Accessibility Reports</h2>
      <hr class="w-full mb-4" />
      <ReportButton
        class="mb-4"
        :game="gameID"
        @submit="
          refreshNuxtData([
            `reports-${gameID}`,
            `featureStats-${gameID}`,
            `score-${gameID}`,
          ])
        "
      />
      <ReportCardGrid :reports="reports" />
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
const reports = ref<Report[]>([]);
const featureStats = ref<FeatureStat[] | null>(null);
const score = ref<number | null>(null);

// Fetch game data
const { data: gameResponse, error: gameError } = await useFetch<Game>(
  "http://localhost:8080/games",
  {
    method: "GET",
    query: {
      id: gameID,
    },
  }
);

if (gameError.value) {
  console.error(gameError.value);
} else if (gameResponse.value) {
  game.value = gameResponse.value;
  const platformPairs = [
    [48, 167], // PS4 and PS5
    [49, 168], // Xbox One and Xbox Series X
  ];
  platformPairs.forEach(([oldPlatform, newPlatform]) => {
    if (
      game.value!.platforms?.includes(oldPlatform) &&
      game.value!.platforms?.includes(newPlatform)
    ) {
      game.value!.platforms = game.value!.platforms.filter(
        (platform) => platform !== oldPlatform
      );
    }
  });
}

// Fetch reports data
const { data: reportsResponse, error: reportError } =
  await useLazyAsyncData<string>(`reports-${gameID}`, () =>
    $fetch(`http://localhost:8080/reports/cards/${gameID}`, {
      method: "GET",
    })
  );

if (reportError.value) {
  console.error(reportError.value);
} else if (reportsResponse.value) {
  reports.value = JSON.parse(reportsResponse.value);
}

// Fetch feature stats data
const { data: featureStatsResponse, error: featureStatsError } =
  await useLazyAsyncData<string>(`featureStats-${gameID}`, () =>
    $fetch(`http://localhost:8080/reports/features/${gameID}`, {
      method: "GET",
    })
  );

if (featureStatsError.value) {
  console.log("error:", featureStatsError.value);
  featureStats.value = null;
} else if (featureStatsResponse.value) {
  console.log(featureStatsResponse.value);
  featureStats.value = JSON.parse(featureStatsResponse.value);
}

// Fetch score data
const { data: scoreResponse, error: scoreError } =
  await useLazyAsyncData<string>(`score-${gameID}`, () =>
    $fetch(`http://localhost:8080/reports/score/${gameID}`, {
      method: "GET",
    })
  );

if (scoreError.value) {
  console.error(scoreError.value);
} else if (scoreResponse.value) {
  score.value = parseFloat(scoreResponse.value);
}

// Watch for changes in responses
watch([reportsResponse, featureStatsResponse, scoreResponse], () => {
  if (reportsResponse.value) {
    reports.value = JSON.parse(reportsResponse.value);
  }
  if (featureStatsResponse.value) {
    featureStats.value = JSON.parse(featureStatsResponse.value);
  }
  if (scoreResponse.value) {
    score.value = parseFloat(scoreResponse.value);
  }
});

// Compute platforms
const platforms = computed(() =>
  platformDefinitions.filter((platform) =>
    game.value?.platforms?.includes(platform.id)
  )
);
</script>

<style></style>
