<template>
  <Card class="w-full h-full overflow-hidden shadow-lg dark">
    <NuxtLink :to="`/games/${game.game_id}`" class="flex flex-col h-full">
      <div class="relative pb-[133%]">
        <NuxtImg
          :src="game.cover_art"
          class="absolute inset-0 w-full h-full object-cover"
        />
      </div>
      <CardFooter
        class="flex-grow flex items-center justify-center p-4 bg-opacity-90 bg-gray-800"
      >
        <p class="text-lg font-bold text-white truncate">{{ game.name }}</p>
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
