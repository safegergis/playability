<template>
  <div
    class="flex items-center justify-center shadow-xl size-24 rounded-sm bg-zinc-700 p-2"
  >
    <h1 class="text-5xl font-bold">{{ score }}</h1>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  game: number;
}>();
const score = ref<number>();
const { data, error } = await useFetch<string>(
  `http://localhost:8080/reports/score/${props.game}`,
  {
    method: "GET",
  }
);
if (error.value) {
  console.error(error.value);
}
if (data.value) {
  score.value = parseFloat(data.value);
}
</script>

<style></style>
