<template>
  <Card class="w-full flex bg-primary text-foreground">
    <div
      class="flex-shrink-0 w-24 bg-zinc-600 flex flex-col items-center justify-center p-2 rounded-l-lg"
    >
      <span class="text-sm font-semibold text-center">Accessibility Score</span>
      <span class="text-2xl font-bold">{{ data?.score }}/10</span>
    </div>
    <div class="flex-grow">
      <CardHeader>
        <CardTitle>{{ data?.username }}</CardTitle>
      </CardHeader>
      <CardContent>
        <p>{{ data?.text }}</p>
      </CardContent>
    </div>
  </Card>
</template>

<script lang="ts" setup>
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

const props = defineProps<{
  reviewId: string;
}>();

/* const { data, error } = useFetch(
  `http://localhost:3000/api/reviews/${props.reviewId}`,
  {
    key: props.reviewId,
  }
  
);
*/
const data = {
  score: 8,
  username: "John Doe",
  text: "This is a test review",
};
</script>

<style scoped>
/* Add any component-specific styles here */
</style>
