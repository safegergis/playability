<template>
    <div class="flex flex-col items-center justify-center min-h-screen p-4 dark">
        <!-- Header -->
        <Transition
            enter-active-class="transition-all duration-300 ease-out"
            enter-from-class="opacity-0 -translate-y-4"
            enter-to-class="opacity-100 translate-y-0"
        >
            <div v-if="searchQuery" class="mb-8 text-center">
                <h1 class="text-3xl font-bold mb-2">
                    Search Results
                </h1>
                <p class="text-muted-foreground">
                    for "<span class="font-semibold text-foreground">{{ searchQuery }}</span>"
                </p>
                <!-- Results Count -->
                <Transition
                    enter-active-class="transition-all duration-200 ease-out delay-100"
                    enter-from-class="opacity-0 translate-y-2"
                    enter-to-class="opacity-100 translate-y-0"
                >
                    <p
                        v-if="!isPending && searchResults && searchResults.length > 0"
                        class="mt-2 text-sm text-muted-foreground"
                        role="status"
                        aria-live="polite"
                    >
                        Found {{ searchResults.length }} {{ searchResults.length === 1 ? 'game' : 'games' }}
                    </p>
                </Transition>
            </div>
        </Transition>

        <!-- Loading State -->
        <Transition
            enter-active-class="transition-all duration-200 ease-out"
            enter-from-class="opacity-0 scale-95"
            enter-to-class="opacity-100 scale-100"
        >
            <div
                v-if="isPending"
                class="w-full max-w-2xl space-y-4"
                role="status"
                aria-live="polite"
                aria-label="Loading search results"
            >
                <Card
                    v-for="i in 5"
                    :key="i"
                    class="dark rounded-lg p-4 shadow-md"
                >
                    <div class="h-6 w-3/4 bg-muted animate-pulse rounded" />
                </Card>
            </div>
        </Transition>

        <!-- Search Results -->
        <TransitionGroup
            v-if="!isPending && searchResults && searchResults.length > 0"
            enter-active-class="transition-all duration-300 ease-out"
            enter-from-class="opacity-0 translate-y-4"
            enter-to-class="opacity-100 translate-y-0"
            leave-active-class="transition-all duration-200 ease-in"
            leave-from-class="opacity-100 translate-y-0"
            leave-to-class="opacity-0 -translate-y-4"
            tag="ul"
            class="w-full max-w-2xl space-y-4"
            role="list"
            aria-label="Search results"
        >
            <li
                v-for="(result, index) in searchResults"
                :key="result.id"
                :style="{ transitionDelay: `${index * 50}ms` }"
            >
                <NuxtLink
                    :to="`/games/${result.id}`"
                    class="block focus:outline-none"
                    :aria-label="`View details for ${result.name}`"
                >
                    <Card
                        class="dark rounded-lg p-6 shadow-md transition-all duration-200 hover:-translate-y-1 hover:shadow-lg hover:border-primary/50 focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 cursor-pointer group"
                    >
                        <h2 class="text-xl font-semibold transition-colors duration-200 group-hover:text-primary">
                            {{ result.name }}
                        </h2>
                    </Card>
                </NuxtLink>
            </li>
        </TransitionGroup>

        <!-- No Results State -->
        <Transition
            enter-active-class="transition-all duration-300 ease-out"
            enter-from-class="opacity-0 scale-95"
            enter-to-class="opacity-100 scale-100"
        >
            <Card
                v-if="!isPending && searchQuery && searchResults && searchResults.length === 0"
                class="dark rounded-lg p-12 shadow-md text-center max-w-md"
            >
                <div class="flex flex-col items-center gap-4">
                    <div
                        class="flex h-16 w-16 items-center justify-center rounded-full bg-muted"
                        aria-hidden="true"
                    >
                        <Icon name="lucide:search-x" class="h-8 w-8 text-muted-foreground" />
                    </div>
                    <div class="space-y-2">
                        <h2 class="text-xl font-semibold">No Results Found</h2>
                        <p class="text-muted-foreground">
                            No games found for "<span class="font-semibold text-foreground">{{ searchQuery }}</span>"
                        </p>
                        <p class="text-sm text-muted-foreground">
                            Try adjusting your search terms
                        </p>
                    </div>
                </div>
            </Card>
        </Transition>
    </div>
</template>

<script lang="ts" setup>
const route = useRoute();
const searchQuery = computed(() => route.query.s as string);

const config = useRuntimeConfig();

const { data: searchResults, status } = await useFetch<SearchResult[]>(
    '/search',
    {
        baseURL: import.meta.server ? config.apiUrl : config.public.apiUrl,
        query: { search: searchQuery },
        default: () => [],
        watch: [searchQuery]
    }
);

const onSearchResultClick = (id: number) => {
    navigateTo(`/games/${id}`);
};
</script>

<style></style>
