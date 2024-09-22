<template>
  <Card class="dark w-full flex">
    <div
      class="flex-shrink-0 w-24 bg-zinc-600 flex flex-col items-center justify-center p-2 rounded-l-lg"
    >
      <span class="text-sm font-semibold text-center">Accessibility Score</span>
      <span class="text-2xl font-bold">{{ report?.score }}/10</span>
    </div>
    <div class="flex-grow">
      <CardHeader>
        <CardTitle>{{ report?.username }}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>{{ report?.report }}</p>
      </CardContent>
    </div>
  </Card>
</template>

<script lang="ts" setup>
interface Report {
  id: string;
  game_id: string;
  user_id: string;
  score: number;
  report: string;
  username?: string;
}
const props = defineProps<{
  report: Report;
}>();
const { data, error } = await useFetch<string>(
  `http://localhost:8080/user/${props.report.user_id}`,
  {
    method: "GET",
  }
);
if (error) {
  console.error(error);
}
if (data.value) {
  console.log("Data: ", data.value);
  const user = JSON.parse(data.value);
  props.report.username = user.username;
}
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
