<template>
  <div class="min-h-screen p-4 dark">
    <div class="container mx-auto max-w-6xl">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-4xl font-bold mb-2">My Accessibility Reports</h1>
        <p class="text-muted-foreground">
          View and manage all the accessibility reports you've submitted
        </p>
      </div>

      <!-- Loading State -->
      <div
        v-if="isLoading"
        class="space-y-4"
        role="status"
        aria-live="polite"
        aria-label="Loading reports"
      >
        <Card
          v-for="i in 3"
          :key="i"
          class="dark p-6"
        >
          <div class="space-y-4">
            <div class="flex items-start justify-between gap-4">
              <div class="flex items-start gap-4 flex-1">
                <!-- Cover Art Skeleton -->
                <div class="h-20 w-16 bg-muted animate-pulse rounded-lg flex-shrink-0" />
                <!-- Info Skeleton -->
                <div class="space-y-2 flex-1">
                  <div class="h-6 w-3/4 bg-muted animate-pulse rounded" />
                  <div class="h-4 w-1/2 bg-muted animate-pulse rounded" />
                </div>
              </div>
              <!-- Score Badge Skeleton -->
              <div class="h-12 w-12 bg-muted animate-pulse rounded-full flex-shrink-0" />
            </div>
            <div class="h-20 w-full bg-muted animate-pulse rounded" />
          </div>
        </Card>
      </div>

      <!-- Empty State -->
      <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
      >
        <Card
          v-if="!isLoading && reports.length === 0"
          class="dark p-12 text-center"
        >
          <div class="flex flex-col items-center gap-4">
            <div
              class="flex h-20 w-20 items-center justify-center rounded-full bg-muted transition-all duration-200"
              aria-hidden="true"
            >
              <Icon name="lucide:file-text" class="h-10 w-10 text-muted-foreground" />
            </div>
            <div class="space-y-2">
              <h2 class="text-2xl font-semibold">No Reports Yet</h2>
              <p class="text-muted-foreground max-w-md">
                You haven't submitted any accessibility reports. Help the community by sharing your gaming experiences!
              </p>
            </div>
            <NuxtLink to="/search" class="mt-4">
              <Button
                size="lg"
                class="transition-all duration-200 hover:scale-105 hover:shadow-md active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
              >
                <Icon name="lucide:search" class="mr-2 h-4 w-4" aria-hidden="true" />
                Find Games to Review
              </Button>
            </NuxtLink>
          </div>
        </Card>
      </Transition>

      <!-- Reports List -->
      <TransitionGroup
        v-if="!isLoading && reports.length > 0"
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 -translate-y-4"
        class="space-y-4"
        tag="div"
      >
        <Card
          v-for="(report, index) in reports"
          :key="report.id"
          class="dark p-6 transition-all duration-300 hover:-translate-y-1 hover:shadow-lg"
          :style="{ transitionDelay: `${index * 50}ms` }"
        >
          <div class="space-y-4">
            <!-- Report Header -->
            <div class="flex items-start justify-between gap-4">
              <div class="flex items-start gap-4 flex-1 min-w-0">
                <!-- Game Cover Art -->
                <NuxtLink
                  v-if="report.cover_art"
                  :to="`/games?id=${report.game_id}`"
                  class="group flex-shrink-0"
                >
                  <div class="relative h-20 w-16 overflow-hidden rounded-lg transition-all duration-200 group-hover:scale-105 group-hover:shadow-lg">
                    <img
                      :src="report.cover_art"
                      :alt="`${report.game_name} cover art`"
                      class="h-full w-full object-cover"
                      loading="lazy"
                    />
                  </div>
                </NuxtLink>

                <!-- Game Info -->
                <div class="flex-1 min-w-0">
                  <NuxtLink
                    :to="`/games?id=${report.game_id}`"
                    class="group inline-block"
                  >
                    <h3
                      class="text-xl font-semibold transition-colors duration-200 group-hover:text-primary group-focus-visible:text-primary truncate"
                    >
                      {{ report.game_name }}
                    </h3>
                  </NuxtLink>
                  <div class="flex items-center gap-3 mt-2 text-sm text-muted-foreground flex-wrap">
                    <div class="flex items-center gap-1.5">
                      <Icon name="lucide:calendar" class="h-4 w-4" aria-hidden="true" />
                      <time :datetime="report.created_at">
                        {{ formatDate(report.created_at) }}
                      </time>
                    </div>
                    <div class="flex items-center gap-1.5">
                      <Icon name="lucide:gamepad-2" class="h-4 w-4" aria-hidden="true" />
                      <span>{{ getPlatformName(report.platform) }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Score Badge -->
              <div
                class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-full transition-all duration-200 hover:scale-110"
                :class="getScoreColorClass(report.score)"
                role="status"
                :aria-label="`Accessibility score: ${report.score} out of 10`"
              >
                <span class="text-lg font-bold">{{ report.score }}</span>
              </div>
            </div>

            <!-- Report Text -->
            <div
              v-if="report.report"
              class="rounded-lg bg-muted/20 p-4 transition-colors duration-200 hover:bg-muted/30"
            >
              <p class="text-sm leading-relaxed">{{ report.report }}</p>
            </div>

            <!-- Actions -->
            <div class="flex gap-2 pt-2">
              <NuxtLink
                :to="`/games?id=${report.game_id}`"
                class="flex-1"
              >
                <Button
                  variant="outline"
                  class="w-full transition-all duration-200 hover:scale-[1.02] active:scale-[0.98] focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
                >
                  <Icon name="lucide:arrow-right" class="mr-2 h-4 w-4" aria-hidden="true" />
                  View Game
                </Button>
              </NuxtLink>
            </div>
          </div>
        </Card>
      </TransitionGroup>

      <!-- Stats Summary (if reports exist) -->
      <Transition
        enter-active-class="transition-all duration-300 ease-out delay-300"
        enter-from-class="opacity-0 translate-y-4"
        enter-to-class="opacity-100 translate-y-0"
      >
        <Card
          v-if="!isLoading && reports.length > 0"
          class="dark mt-8 p-6"
        >
          <h2 class="text-lg font-semibold mb-4">Your Impact</h2>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="text-center p-4 rounded-lg bg-muted/20">
              <div class="text-3xl font-bold text-primary mb-1">{{ reports.length }}</div>
              <div class="text-sm text-muted-foreground">
                Total {{ reports.length === 1 ? 'Report' : 'Reports' }}
              </div>
            </div>
            <div class="text-center p-4 rounded-lg bg-muted/20">
              <div class="text-3xl font-bold text-primary mb-1">{{ averageScore }}</div>
              <div class="text-sm text-muted-foreground">Average Score</div>
            </div>
            <div class="text-center p-4 rounded-lg bg-muted/20">
              <div class="text-3xl font-bold text-primary mb-1">{{ uniqueGames }}</div>
              <div class="text-sm text-muted-foreground">
                Unique {{ uniqueGames === 1 ? 'Game' : 'Games' }}
              </div>
            </div>
          </div>
        </Card>
      </Transition>
    </div>
  </div>
</template>

<script lang="ts" setup>

// Use SEO meta tags
useHead({
  title: 'My Reports - Playability',
  meta: [
    {
      name: 'description',
      content: 'View all your accessibility reports and contributions to the gaming community',
    },
  ],
});

interface Report {
  id: number;
  created_at: string;
  game_id: number;
  game_name: string;
  cover_art: string;
  user_id: number;
  platform: number;
  score: number;
  report: string;
}

const authStore = useAuthStore();
const isLoading = ref(true);
const reports = ref<Report[]>([]);

// Fetch user reports on mount
onMounted(async () => {
  if (!authStore.loggedIn) {
    navigateTo('/login');
    return;
  }

  try {
    const data = await $fetch<Report[]>('/api/user/my-reports');
    reports.value = data || [];
  } catch (error) {
    console.error('Error fetching reports:', error);
    // Could add a toast notification here
  } finally {
    isLoading.value = false;
  }
});

// Helper: Format date
const formatDate = (dateString: string) => {
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date);
};

// Helper: Get platform name
const getPlatformName = (platform: number) => {
  const platforms: Record<number, string> = {
    49: 'PlayStation',
    169: 'Xbox',
    6: 'PC',
    130: 'Nintendo',
  };
  return platforms[platform] || 'Unknown';
};

// Helper: Get score color class
const getScoreColorClass = (score: number) => {
  if (score >= 8) return 'bg-green-500/20 text-green-500';
  if (score >= 5) return 'bg-yellow-500/20 text-yellow-500';
  return 'bg-red-500/20 text-red-500';
};

// Computed: Average score
const averageScore = computed(() => {
  if (reports.value.length === 0) return 0;
  const sum = reports.value.reduce((acc, report) => acc + report.score, 0);
  return (sum / reports.value.length).toFixed(1);
});

// Computed: Unique games count
const uniqueGames = computed(() => {
  const gameIds = new Set(reports.value.map(r => r.game_id));
  return gameIds.size;
});
</script>
