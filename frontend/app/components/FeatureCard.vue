<template>
  <Card class="dark p-6 rounded-xl animate-fade-in">
    <h2 class="text-2xl font-bold mb-6 flex items-center gap-2">
      <Icon name="lucide:accessibility" class="w-6 h-6 text-primary" aria-hidden="true" />
      Essential Accessibility Features
    </h2>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div
        v-for="feature in featureData"
        :key="feature.name"
        class="group bg-card/50 border border-border p-5 rounded-xl transition-all duration-300 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50"
      >
        <h3 class="text-lg font-semibold mb-3 flex items-center gap-2">
          <Icon :name="feature.icon" class="w-5 h-5 text-xl text-primary transition-transform duration-300 group-hover:scale-110" aria-hidden="true" />
          {{ feature.name }}
        </h3>
        <p class="text-muted-foreground font-medium mb-1">
          {{ feature.consensusText }}
        </p>
        <p
          v-if="feature.showSecondaryConsensus"
          class="text-muted-foreground/80 text-sm"
        >
          {{ feature.secondaryConsensusText }}
        </p>
      </div>
    </div>
  </Card>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
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

const featureData = computed(() => {
  const features = [
    { name: 'Colorblind Mode', icon: 'lucide:palette', gameProp: 'color_blind' },
    { name: 'Closed Captions', icon: 'lucide:closed-caption', gameProp: 'closed_captions' },
    { name: 'Controller Support', icon: 'lucide:gamepad-2', gameProp: 'full_controller_support' },
    { name: 'Full Controller Remapping', icon: 'lucide:settings', gameProp: 'controller_remapping' }
  ];

  return features.map((feature, index) => {
    const stat = props.featureStats ? props.featureStats[index] : null;
    
    const consensusValue = stat ? stat.consensus : (props.game[feature.gameProp as keyof Game] as string || 'unknown');
    const consensusText = parseConsensus(consensusValue);

    const secondaryConsensusValue = stat ? stat.secondary_consensus : '';
    const showSecondaryConsensus = ifSecondaryConsensus(secondaryConsensusValue);
    const secondaryConsensusText = stat && showSecondaryConsensus ? parseSecondaryConsensus(stat) : '';

    return {
      ...feature,
      consensusText,
      showSecondaryConsensus,
      secondaryConsensusText
    };
  });
});
</script>

<style></style>
