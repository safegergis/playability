<template>
    <main class="min-h-screen">
        <!-- Game Header Section -->
        <div class="container mx-auto my-auto px-4 py-8 lg:py-12">
            <div v-if="game" class="space-y-8">
                <!-- Game Info Section -->
                <section class="flex flex-col lg:flex-row gap-8 lg:gap-12 animate-fade-in">
                    <!-- Cover Art -->
                    <div class="lg:1-1/3">
                        <div class="flex justify-center lg:justify-start">
                            <div
                                class="group relative overflow-hidden rounded-2xl shadow-2xl transition-all duration-500 hover:scale-105 hover:shadow-primary/20">
                                <NuxtImg :src="game.cover_art" :alt="`${game.name} cover art`" width="400" height="533"
                                    class="rounded-2xl transition-transform duration-500 group-hover:scale-110"
                                    loading="eager" />
                                <div
                                    class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100 rounded-2xl" />
                            </div>
                        </div>
                    </div>

                    <!-- Game Details -->
                    <div class="lg:w-2/3 space-y-6">
                        <!-- Title and Score -->
                        <div class="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                            <h1
                                class="text-4xl lg:text-5xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text">
                                {{ game.name }}
                            </h1>
                            <AccesibilityScoreBadge :score="score" />
                        </div>

                        <!-- Platforms -->
                        <div v-if="game.platforms && platforms.length > 0" class="dark space-y-3">
                            <h2 class="text-lg font-semibold text-muted-foreground flex items-center gap-2">
                                <Icon name="lucide:monitor" class="w-5 h-5" aria-hidden="true" />
                                Available Platforms
                            </h2>
                            <ul class="flex flex-wrap gap-3" role="list" aria-label="Game platforms">
                                <li v-for="platform in platforms" :key="platform.id"
                                    class="group bg-card border border-border px-4 py-2 rounded-lg transition-all duration-200 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50">
                                    <Icon :name="platform.icon"
                                        class="w-6 h-6 transition-transform duration-200 group-hover:scale-110"
                                        :aria-label="platform.name" />
                                </li>
                            </ul>
                        </div>

                        <!-- Summary -->
                        <div v-if="game.summary" class="space-y-3">
                            <h2 class="text-lg font-semibold text-muted-foreground flex items-center gap-2">
                                <Icon name="lucide:info" class="w-5 h-5" aria-hidden="true" />
                                About
                            </h2>
                            <p class="text-muted-foreground leading-relaxed text-base lg:text-lg">
                                {{ game.summary }}
                            </p>
                        </div>

                        <!-- Accessibility Features -->
                        <FeatureCard :feature-stats="featureStats" :game="game" />
                    </div>
                </section>
            </div>

            <!-- Loading State -->
            <div v-else class="flex flex-col items-center justify-center min-h-[50vh] space-y-4" role="status"
                aria-live="polite">
                <Icon name="lucide:loader-2" class="w-16 h-16 animate-spin text-primary" aria-hidden="true" />
                <p class="text-xl text-muted-foreground">Loading game details...</p>
                <span class="sr-only">Loading game information, please wait</span>
            </div>
        </div>

        <!-- Reports Section -->
        <div class="bg-neutral-950 border-y border-border py-8 lg:py-12 mt-8">
            <div class="container mx-auto px-4 space-y-6">
                <!-- Section Header -->
                <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 animate-slide-up">
                    <div>
                        <h2 class="text-3xl lg:text-4xl font-bold flex items-center gap-3">
                            <Icon name="lucide:messages-square" class="w-8 h-8 text-neutral-500" aria-hidden="true" />
                            Accessibility Reports
                        </h2>
                        <p class="text-muted-foreground mt-2">
                            Real experiences from players like you
                        </p>
                    </div>
                    <ReportButton :game="gameID" @submit="
                        refreshNuxtData([
                            `reports-${gameID}`,
                            `featureStats-${gameID}`,
                            `score-${gameID}`,
                        ])
                        " />
                </div>

                <!-- Reports Grid -->
                <div class="animate-slide-up" style="animation-delay: 150ms">
                    <ReportCardGrid :reports="reports" />
                </div>
            </div>
        </div>
    </main>
</template>

<script lang="ts" setup>
const platformDefinitions = [
    { id: 48, name: "Playstation 4", icon: "mdi:sony-playstation" },
    { id: 49, name: "Xbox One", icon: "mdi:microsoft-xbox" },
    { id: 6, name: "PC", icon: "mdi:steam" },
    { id: 130, name: "Nintendo Switch", icon: "mdi:nintendo-switch" },
    { id: 167, name: "Playstation 5", icon: "mdi:sony-playstation" },
    { id: 168, name: "Xbox Series X", icon: "mdi:microsoft-xbox" },
];

const route = useRoute();
const gameID = computed(() => Number.parseInt(route.params.id as string));

const config = useRuntimeConfig();
const baseURL = import.meta.server ? config.apiUrl : config.public.apiUrl;

// Fetch game data
const { data: game } = await useFetch<Game>(
    '/games',
    {
        baseURL,
        query: { id: gameID },
        watch: [gameID],
        transform: (data) => {
            if (data) {
                const platformPairs = [
                    [48, 167], // PS4 and PS5
                    [49, 168], // Xbox One and Xbox Series X
                ];
                platformPairs.forEach(([oldPlatform, newPlatform]) => {
                    if (
                        data.platforms?.includes(oldPlatform) &&
                        data.platforms?.includes(newPlatform)
                    ) {
                        data.platforms = data.platforms.filter(
                            (platform) => platform !== oldPlatform
                        );
                    }
                });
            }
            return data;
        }
    }
);

// Fetch reports data
const { data: reports } = await useFetch<Report[]>(
    () => `/reports/cards/${gameID.value}`,
    {
        baseURL,
        key: `reports-${gameID.value}`,
        default: () => [],
        watch: [gameID]
    }
);

// Fetch feature stats data
const { data: featureStats } = await useFetch<FeatureStat[] | null>(
    () => `/reports/features/${gameID.value}`,
    {
        baseURL,
        key: `featureStats-${gameID.value}`,
        default: () => null,
        watch: [gameID]
    }
);

// Fetch score data
const { data: score } = await useFetch<number | null>(
    () => `/reports/score/${gameID.value}`,
    {
        baseURL,
        key: `score-${gameID.value}`,
        transform: (data) => data ? parseFloat(data as any) : null,
        default: () => null,
        watch: [gameID]
    }
);

// Compute platforms
const platforms = computed(() =>
    platformDefinitions.filter((platform) =>
        game.value?.platforms?.includes(platform.id)
    )
);

// Set page metadata for SEO and accessibility
useHead(() => ({
    title: game.value ? `${game.value.name} - Accessibility Information | Playability` : 'Game Details | Playability',
    meta: [
        {
            name: 'description',
            content: game.value?.summary
                ? `${game.value.summary.substring(0, 155)}...`
                : 'View detailed accessibility information and user reports for this game.'
        }
    ]
}));

</script>

<style></style>
