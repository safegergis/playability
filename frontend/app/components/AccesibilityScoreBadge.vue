<template>
  <div
    v-if="score"
    class="group relative flex flex-col items-center justify-center shadow-2xl size-32 p-3 rounded-2xl transition-all duration-300 hover:scale-105 hover:shadow-2xl animate-scale-in"
    :class="badgeClasses"
    role="status"
    :aria-label="`Accessibility score: ${scoreLabel} - ${score} out of 10`"
  >
    <!-- Glow effect -->
    <div
      class="absolute inset-0 rounded-2xl blur-xl opacity-40 transition-opacity duration-300 group-hover:opacity-60"
      :class="glowClasses"
    />

    <!-- Inner content -->
    <div class="relative z-10 flex flex-col items-center">
      <span class="text-[10px] font-bold uppercase tracking-wider mb-1" :class="labelClasses">
        {{ scoreLabel }}
      </span>
      <span class="text-5xl font-bold transition-transform duration-300 group-hover:scale-110" :class="textClasses">
        {{ score }}
      </span>
      <span class="text-xs mt-1 opacity-90" :class="textClasses">/ 10</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{
  score: number | null;
}>();

// Color coding based on score
const scoreLevel = computed(() => {
  if (!props.score) return 'none';
  if (props.score >= 9) return 'excellent';
  if (props.score >= 7) return 'good';
  if (props.score >= 4) return 'fair';
  return 'poor';
});

const scoreLabel = computed(() => {
  switch (scoreLevel.value) {
    case 'excellent': return 'Excellent';
    case 'good': return 'Good';
    case 'fair': return 'Fair';
    case 'poor': return 'Poor';
    default: return 'Score';
  }
});

const badgeClasses = computed(() => {
  switch (scoreLevel.value) {
    case 'excellent':
      return 'bg-gradient-to-br from-emerald-600 to-emerald-800 border-2 border-emerald-400/30';
    case 'good':
      return 'bg-gradient-to-br from-green-600 to-green-800 border-2 border-green-400/30';
    case 'fair':
      return 'bg-gradient-to-br from-amber-600 to-amber-800 border-2 border-amber-400/30';
    case 'poor':
      return 'bg-gradient-to-br from-rose-600 to-rose-800 border-2 border-rose-400/30';
    default:
      return 'bg-gradient-to-br from-primary to-primary/80';
  }
});

const glowClasses = computed(() => {
  switch (scoreLevel.value) {
    case 'excellent':
      return 'bg-emerald-400';
    case 'good':
      return 'bg-green-400';
    case 'fair':
      return 'bg-amber-400';
    case 'poor':
      return 'bg-rose-400';
    default:
      return 'bg-primary';
  }
});

const textClasses = computed(() => {
  return 'text-white font-bold';
});

const labelClasses = computed(() => {
  switch (scoreLevel.value) {
    case 'excellent':
      return 'text-emerald-200';
    case 'good':
      return 'text-green-200';
    case 'fair':
      return 'text-amber-200';
    case 'poor':
      return 'text-rose-200';
    default:
      return 'text-primary-foreground/80';
  }
});
</script>

<style></style>
