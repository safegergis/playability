<template>
  <Card
    v-if="
      featureStats &&
      colorBlind &&
      closedCaptions &&
      controllerSupport &&
      controllerRemapping
    "
    class="dark p-4 rounded-lg"
  >
    <h2 class="text-xl font-semibold mb-4">Essential Accessibility Features</h2>
    <div class="grid grid-cols-2 md:grid-cols-2 gap-4">
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Colorblind Mode</h3>
        <p class="text-gray-300">
          {{ parseConsensus(colorBlind.consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(colorBlind.secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(colorBlind) }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Closed Captions</h3>
        <p class="text-gray-300">
          {{ parseConsensus(closedCaptions.consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(closedCaptions.secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(closedCaptions) }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Controller Support</h3>
        <p class="text-gray-300">
          {{ parseConsensus(controllerSupport.consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(controllerSupport.secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(controllerSupport) }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Full Controller Remapping</h3>
        <p class="text-gray-300">
          {{ parseConsensus(controllerRemapping.consensus) }}
        </p>

        <p
          v-if="ifSecondaryConsensus(controllerRemapping.secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(controllerRemapping) }}
        </p>
      </div>
    </div>
  </Card>

  <Card v-else class="dark p-4 rounded-lg:">
    <h2 class="text-xl font-semibold mb-4">Essential Accessibility Features</h2>

    <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Colorblind Mode</h3>
        <p class="text-gray-300">
          {{ parseConsensus(game.color_blind || "unknown") }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Closed Captions</h3>
        <p class="text-gray-300">
          {{ parseConsensus(game.closed_captions || "unknown") }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Controller Support</h3>
        <p class="text-gray-300">
          {{ parseConsensus(game.full_controller_support || "unknown") }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Full Controller Remapping</h3>
        <p class="text-gray-300">
          {{ parseConsensus(game.controller_remapping || "unknown") }}
        </p>
      </div>
    </div>
  </Card>
</template>

<script setup lang="ts">
const props = defineProps<{
  gameid: number;
  game: Game;
}>();

const fallback = ref(false);
const featureStats = ref<FeatureStat[] | null>(null);

const colorBlind = ref<FeatureStat | null>(null);
const closedCaptions = ref<FeatureStat | null>(null);
const controllerSupport = ref<FeatureStat | null>(null);
const controllerRemapping = ref<FeatureStat | null>(null);
const { data, error } = await useFetch<string>(
  `http://localhost:8080/reports/features/${props.gameid}`,
  {
    method: "GET",
  }
);
if (error.value) {
  fallback.value = true;
  console.log("error:", error.value);
} else if (data.value) {
  console.log(data.value);
  featureStats.value = JSON.parse(data.value);
  if (featureStats.value) {
    colorBlind.value = featureStats.value[0];
    closedCaptions.value = featureStats.value[1];
    controllerSupport.value = featureStats.value[2];
    controllerRemapping.value = featureStats.value[3];
  }
}
const consensusMap = {
  unknown: "Data not available",
  "": "Data not available",
  limited: "Limited implementation",
  true: "Supported",
  false: "Not supported",
};

const parseConsensus = (resp: string) => {
  return consensusMap[resp] || "Data not available";
};
const ifSecondaryConsensus = (resp: string): boolean => {
  return ["true", "limited"].includes(resp);
};

const parseSecondaryConsensus = (feature: FeatureStat) => {
  const consensusMap = {
    true: () =>
      `${(feature.true_percentage * 100).toFixed(
        0
      )}% of users reported full support`,
    limited: () =>
      `${(feature.limited_percentage * 100).toFixed(
        0
      )}% of users reported partial support`,
    false: () => "",
    "": () => "",
  };

  return (consensusMap[feature.secondary_consensus] || (() => ""))();
};
</script>

<style></style>
