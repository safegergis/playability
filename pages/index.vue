<template>
  <div class="dark flex flex-col items-center justify-center min-h-screen">
    <h1 class="text-4xl font-bold mb-4">Playability</h1>
    <div class="flex items-center">
      <Input
        v-model="searchQuery"
        class="w-64 mr-2 text-muted-foreground"
        placeholder="Search games"
      />
      <Button
        variant="outline"
        type="icon"
        class="text-muted-foreground"
        @click="onSearch"
      >
        <Icon name="formkit:search" />
      </Button>
    </div>
  </div>
</template>

<script setup lang="ts">
const searchQuery = ref("");

const onSearch = () => {
  navigateTo(`/search/${searchQuery.value}`, {
    external: false,
  });
};
</script>

<style></style>
