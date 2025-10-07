<template>
    <!-- Success State: Carousel -->
    <div v-if="status === 'success'" class="animate-fade-in">
        <Carousel :opts="{ loop: true }" :plugins="[
            Autoplay({
                delay: 2000,
            }),
        ]" aria-label="Featured games carousel">
            <CarouselContent>
                <CarouselItem v-for="game in games" :key="game.id" class="basis-1/2 md:basis-1/4 lg:basis-1/5"
                    role="group" :aria-label="`Featured game ${game.name}`">
                    <FeaturedGame :game="game" />
                </CarouselItem>
            </CarouselContent>
        </Carousel>
    </div>

    <!-- Loading State: Elegant Skeleton -->
    <div v-if="status === 'pending'" role="status" aria-live="polite" aria-label="Loading featured games"
        class="space-y-4">
        <div class="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            <div v-for="i in 5" :key="i" class="relative overflow-hidden rounded-lg bg-muted"
                :style="{ animationDelay: `${i * 100}ms` }">
                <!-- Skeleton Card -->
                <div class="pb-[133%] relative">
                    <div
                        class="absolute inset-0 bg-gradient-to-r from-muted via-muted-foreground/10 to-muted animate-shimmer bg-[length:1000px_100%]" />
                </div>
                <div class="h-16 bg-muted/80 animate-pulse" />
            </div>
        </div>
        <span class="sr-only">Loading featured games, please wait</span>
    </div>

    <!-- Error State: Friendly Error Message -->
    <div v-else-if="status === 'error'" role="alert" aria-live="assertive"
        class="flex flex-col items-center justify-center h-40 space-y-4 animate-scale-in">
        <Icon name="lucide:alert-circle" class="w-12 h-12 text-destructive animate-pulse" aria-hidden="true" />
        <p class="text-destructive text-lg font-medium text-center max-w-md">
            Unable to load featured games. Please check your connection and try again.
        </p>
    </div>
</template>

<script lang="ts" setup>
import Autoplay from "embla-carousel-autoplay";
defineProps<{
    games: FeaturedGame[];
    status: "pending" | "success" | "error" | "idle";
}>();
</script>

<style></style>
