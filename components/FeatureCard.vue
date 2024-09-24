<template>
  <Card v-if="featureStats" class="dark p-4 rounded-lg">
    <h2 class="text-xl font-semibold mb-4">Essential Accessibility Features</h2>
    <div class="grid grid-cols-2 md:grid-cols-2 gap-4">
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Colorblind Mode</h3>
        <p class="text-gray-300">
          {{ parseConsensus(featureStats[0].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[0].secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[0]) }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Closed Captions</h3>
        <p class="text-gray-300">
          {{ parseConsensus(featureStats[1].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[1].secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[1]) }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Controller Support</h3>
        <p class="text-gray-300">
          {{ parseConsensus(featureStats[2].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[2].secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[2]) }}
        </p>
      </div>
      <div class="bg-zinc-700 p-4 rounded-lg">
        <h3 class="text-lg font-medium mb-2">Full Controller Remapping</h3>
        <p class="text-gray-300">
          {{ parseConsensus(featureStats[3].consensus) }}
        </p>

        <p
          v-if="ifSecondaryConsensus(featureStats[3].secondary_consensus)"
          class="text-gray-300 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[3]) }}
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
defineProps<{
  featureStats: FeatureStat[] | null;
  game: Game;
}>();

const consensusMap = {
  unknown: "Data not available",
  "": "Data not available",
  limited: "Limited implementation",
  true: "Supported",
  false: "Not supported",
};

const parseConsensus = (resp: string) => {
  return (
    consensusMap[resp as keyof typeof consensusMap] || "Data not available"
  );
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

  return (
    consensusMap[feature.secondary_consensus as keyof typeof consensusMap] ||
    (() => "")
  )();
};
</script>

<style></style>
