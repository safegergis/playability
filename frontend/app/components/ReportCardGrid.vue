<template>
  <!-- Reports Grid -->
  <div
    v-if="reports && reports.length > 0"
    class="grid grid-cols-1 lg:grid-cols-2 gap-6"
    role="list"
    aria-label="User accessibility reports"
  >
    <div
      v-for="(report, index) in reports"
      :key="report.id"
      role="listitem"
      class="animate-slide-up"
      :style="{ animationDelay: `${index * 100}ms` }"
    >
      <ReportCard :report="report" />
    </div>
  </div>

  <!-- Empty State -->
  <div
    v-else
    class="flex flex-col items-center justify-center p-12 rounded-xl border-2 border-dashed border-border/50 bg-card/30 transition-all duration-300 hover:border-primary/30 hover:bg-card/50 animate-fade-in"
    role="status"
    aria-label="No reports available"
  >
    <div class="relative">
      <div class="absolute inset-0 bg-primary/10 blur-3xl rounded-full" />
      <Icon
        name="lucide:clipboard-list"
        class="relative w-20 h-20 text-primary/40 mb-6 transition-transform duration-300 hover:scale-110"
        aria-hidden="true"
      />
    </div>
    <h3 class="text-2xl font-bold text-foreground mb-2">
      No Reports Yet
    </h3>
    <p class="text-muted-foreground text-center max-w-md mb-4">
      No accessibility reports available yet.
    </p>
    <p class="text-primary font-medium flex items-center gap-2">
      <Icon name="lucide:sparkles" class="w-4 h-4" aria-hidden="true" />
      Be the first to share your experience!
    </p>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  reports: Report[];
}>();
</script>

<style></style>
