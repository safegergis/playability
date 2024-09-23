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
const props = defineProps<{
  report: Report;
}>();
const report = ref<Report>(props.report);
const user = ref<User>();
const { data, error } = await useFetch<string>(
  `http://localhost:8080/user/${props.report.user_id}`,
  {
    method: "GET",
  }
);
if (error.value) {
  console.error(error.value);
}
if (data.value) {
  user.value = JSON.parse(data.value);
  if (user.value) {
    report.value.username = user.value.username;
  }
}
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
