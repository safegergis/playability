<template>
  <div
    class="flex items-center justify-center shadow-xl size-24 rounded-sm bg-zinc-700 p-2"
  >
    <h1 class="text-5xl font-bold">{{ score }}</h1>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  score: number;
}>();
</script>

<style></style>
