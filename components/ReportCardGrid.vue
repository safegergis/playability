<template>
  <div
    v-if="reports && reports.length > 0"
    class="grid grid-cols-1 sm:grid-cols-2 gap-4"
  >
    <div v-for="report in reports" :key="report.id">
      <ReportCard :report="report" />
    </div>
  </div>
  <div
    v-else
    class="flex flex-col items-center justify-center p-8 bg-zinc-800 rounded-lg shadow-md"
  >
    <Icon
      name="mdi:clipboard-text-outline"
      class="text-6xl text-gray-400 mb-4"
    />
    <p class="text-xl text-gray-300 text-center">
      No accessibility reports available yet.
    </p>
    <p class="text-gray-400 text-center mt-2">
      Be the first to share your experience!
    </p>
  </div>
</template>

<script lang="ts" setup>
interface Report {
  id: string;
  game_id: string;
  user_id: string;
  score: number;
  report: string;
}
const props = defineProps<{
  game: string;
}>();
const { data, error } = await useFetch<Report[]>(
  `http://localhost:8080/reports/cards/${props.game}`,
  {
    method: "GET",
  }
);
if (error) {
  console.error(error);
}
console.log("Data: ", data.value);
const reports = JSON.parse(data.value as unknown as string);
</script>

<style></style>
