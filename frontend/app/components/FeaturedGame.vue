<template>
  <Card class="group w-full h-full overflow-hidden shadow-lg transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 dark">
    <NuxtLink
      :to="`/games/${game.game_id}`"
      class="flex flex-col h-full focus-visible:outline-none"
      :aria-label="`View details for ${game.name}`"
    >
      <div class="relative pb-[133%] overflow-hidden bg-muted">
        <NuxtImg
          :src="game.cover_art"
          :alt="`${game.name} cover art`"
          class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
          loading="lazy"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
      </div>
      <CardFooter
        class="flex-grow flex items-center justify-center p-4 bg-card/90 backdrop-blur-sm transition-colors duration-300 group-hover:bg-card"
      >
        <p class="text-lg font-bold text-card-foreground truncate transition-transform duration-300 group-hover:scale-105">
          {{ game.name }}
        </p>
      </CardFooter>
    </NuxtLink>
  </Card>
</template>

<script lang="ts" setup>
defineProps<{
  game: FeaturedGame;
}>();
</script>

<style></style>
