<template>
  <Card v-if="featureStats" class="dark p-6 rounded-xl animate-fade-in">
    <h2 class="text-2xl font-bold mb-6 flex items-center gap-2">
      <Icon name="lucide:accessibility" class="w-6 h-6 text-primary" aria-hidden="true" />
      Essential Accessibility Features
    </h2>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <!-- Colorblind Mode -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:palette" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Colorblind Mode
        </h3>
        <p class="text-muted-foreground font-medium mb-1">
          {{ parseConsensus(featureStats[0].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[0].secondary_consensus)"
          class="text-muted-foreground/80 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[0]) }}
        </p>
      </div>

      <!-- Closed Captions -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:closed-captioning" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Closed Captions
        </h3>
        <p class="text-muted-foreground font-medium mb-1">
          {{ parseConsensus(featureStats[1].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[1].secondary_consensus)"
          class="text-muted-foreground/80 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[1]) }}
        </p>
      </div>

      <!-- Controller Support -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:gamepad-2" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Controller Support
        </h3>
        <p class="text-muted-foreground font-medium mb-1">
          {{ parseConsensus(featureStats[2].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[2].secondary_consensus)"
          class="text-muted-foreground/80 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[2]) }}
        </p>
      </div>

      <!-- Full Controller Remapping -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:settings" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Full Controller Remapping
        </h3>
        <p class="text-muted-foreground font-medium mb-1">
          {{ parseConsensus(featureStats[3].consensus) }}
        </p>
        <p
          v-if="ifSecondaryConsensus(featureStats[3].secondary_consensus)"
          class="text-muted-foreground/80 text-sm"
        >
          {{ parseSecondaryConsensus(featureStats[3]) }}
        </p>
      </div>
    </div>
  </Card>

  <Card v-else class="dark p-6 rounded-xl animate-fade-in">
    <h2 class="text-2xl font-bold mb-6 flex items-center gap-2">
      <Icon name="lucide:accessibility" class="w-6 h-6 text-primary" aria-hidden="true" />
      Essential Accessibility Features
    </h2>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Colorblind Mode -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:palette" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Colorblind Mode
        </h3>
        <p class="text-muted-foreground font-medium">
          {{ parseConsensus(game.color_blind || "unknown") }}
        </p>
      </div>

      <!-- Closed Captions -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:closed-captioning" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Closed Captions
        </h3>
        <p class="text-muted-foreground font-medium">
          {{ parseConsensus(game.closed_captions || "unknown") }}
        </p>
      </div>

      <!-- Controller Support -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:gamepad-2" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Controller Support
        </h3>
        <p class="text-muted-foreground font-medium">
          {{ parseConsensus(game.full_controller_support || "unknown") }}
        </p>
      </div>

      <!-- Full Controller Remapping -->
      <div
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon name="lucide:settings" class="w-5 h-5 text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          Full Controller Remapping
        </h3>
        <p class="text-muted-foreground font-medium">
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
