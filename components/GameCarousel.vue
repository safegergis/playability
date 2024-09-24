<template>
  <div v-if="status === 'success'">
    <Carousel
      :opts="{ loop: true }"
      :plugins="[
        Autoplay({
          delay: 2000,
        }),
      ]"
    >
      <CarouselContent>
        <CarouselItem
          v-for="game in games"
          :key="game.id"
          class="basis-1/3 md:basis-1/4 lg:basis-1/5"
        >
          <FeaturedGame :game="game" />
        </CarouselItem>
      </CarouselContent>
    </Carousel>
  </div>
  <div
    v-if="status === 'pending'"
    class="flex justify-center items-center h-40"
  >
    <div class="rounded-full h-20 w-20 bg-gray-800 animate-ping" />
  </div>

  <div
    v-else-if="status === 'error'"
    class="flex justify-center items-center h-40"
  >
    <p class="text-red-500 text-lg">
      Error loading featured games. Please try again later.
    </p>
  </div>
</template>

<script lang="ts" setup>
import Autoplay from "embla-carousel-autoplay";
defineProps<{
  games: FeaturedGame[];
  status: "pending" | "success" | "error";
}>();
</script>

<style></style>
