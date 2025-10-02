<template>
  <nav class="p-4">
    <div class="container mx-auto flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <NuxtLink to="/" class="text-3xl font-bold text-white">
          Playability
        </NuxtLink>
        <form class="relative" @submit.prevent="handleSearch">
          <Input
            v-model="searchQuery"
            type="text"
            placeholder="Search games..."
            class="dark"
          />
          <Button
            type="submit"
            variant="ghost"
            class="absolute right-0 top-0 h-full px-3 dark"
          >
            <Icon name="lucide:search" class="h-4 w-4" />
            <span class="sr-only">Search</span>
          </Button>
        </form>
      </div>

      <div class="flex items-center space-x-4">
        <NavigationMenu class="dark">
          <NavigationMenuList>
            <NavigationMenuItem>
              <NuxtLink to="/about">
                <NavigationMenuLink :class="navigationMenuTriggerStyle()">
                  About
                </NavigationMenuLink>
              </NuxtLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NuxtLink to="/donate">
                <NavigationMenuLink :class="navigationMenuTriggerStyle()">
                  Donate
                </NavigationMenuLink>
              </NuxtLink>
            </NavigationMenuItem>
            <NavigationMenuItem class="ml-8">
              <NuxtLink to="/login">
                <Button variant="link" class="dark">Login</Button>
              </NuxtLink>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <NuxtLink to="/register">
                <Button variant="link" class="dark">Register</Button>
              </NuxtLink>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
      </div>
    </div>
  </nav>
</template>

<script lang="ts" setup>
import { navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";
const searchQuery = ref("");

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    navigateTo({
      path: "/search",
      query: {
        s: searchQuery.value,
      },
    });
  }
};
</script>

<style scoped>
/* Add any additional styles here if needed */
</style>
