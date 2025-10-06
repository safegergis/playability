<template>
    <nav
        class="sticky top-0 dark z-40 w-full border-b border-border/40 bg-background/95 p-4 backdrop-blur supports-[backdrop-filter]:bg-background/60 transition-all duration-300"
        role="navigation"
        aria-label="Main navigation"
    >
        <div class="container mx-auto flex items-center justify-between gap-4">
            <!-- Logo and Brand -->
            <div class="flex items-center gap-6">
                <NuxtLink
                    to="/"
                    class="group flex items-center gap-3 rounded-lg transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:ring-offset-background"
                    aria-label="Playability home"
                >
                    <div class="relative h-10 w-10 overflow-hidden rounded-lg transition-transform duration-200 group-hover:scale-110 group-active:scale-95">
                        <img
                            src="/playability_logo.png"
                            alt="Playability logo"
                            class="h-full w-full object-contain"
                            width="40"
                            height="40"
                        />
                    </div>
                    <span class="text-2xl font-bold text-foreground transition-colors duration-200 group-hover:text-primary">
                        Playability
                    </span>
                </NuxtLink>

                <!-- Search Form -->
                <form
                    class="relative hidden md:block"
                    role="search"
                    aria-label="Search games"
                    @submit.prevent="handleSearch"
                >
                    <div class="relative">
                        <Input
                            id="nav-search-input"
                            v-model="searchQuery"
                            type="search"
                            placeholder="Search games..."
                            class="dark w-64 pr-10 transition-all duration-200 focus-visible:w-80 lg:w-80"
                            aria-label="Search for games"
                            :aria-busy="isSearching"
                        />
                        <Button
                            type="submit"
                            variant="ghost"
                            size="icon"
                            class="absolute right-0 top-0 h-full px-3 dark transition-transform duration-150 hover:scale-110 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-1"
                            :disabled="isSearching || !searchQuery.trim()"
                            :aria-label="isSearching ? 'Searching...' : 'Submit search'"
                        >
                            <Icon
                                v-if="!isSearching"
                                name="lucide:search"
                                class="h-4 w-4 transition-colors"
                                aria-hidden="true"
                            />
                            <Icon
                                v-else
                                name="lucide:loader-2"
                                class="h-4 w-4 animate-spin"
                                aria-hidden="true"
                            />
                            <span class="sr-only">{{ isSearching ? 'Searching...' : 'Search' }}</span>
                        </Button>
                    </div>
                </form>
            </div>

            <!-- Navigation Links -->
            <div class="flex items-center gap-2">
                <NavigationMenu class="dark">
                    <NavigationMenuList>
                        <NavigationMenuItem>
                            <NuxtLink to="/about">
                                <NavigationMenuLink
                                    :class="navigationMenuTriggerStyle()"
                                    class="transition-all duration-200 hover:scale-105 active:scale-95"
                                >
                                    About
                                </NavigationMenuLink>
                            </NuxtLink>
                        </NavigationMenuItem>
                        <NavigationMenuItem>
                            <NuxtLink to="/donate">
                                <NavigationMenuLink
                                    :class="navigationMenuTriggerStyle()"
                                    class="transition-all duration-200 hover:scale-105 active:scale-95"
                                >
                                    Donate
                                </NavigationMenuLink>
                            </NuxtLink>
                        </NavigationMenuItem>

                        <!-- Authentication Links -->
                        <NavigationMenuItem v-if="!authStore.loggedIn" class="ml-4">
                            <NuxtLink to="/login">
                                <Button
                                    variant="link"
                                    class="dark transition-all duration-200 hover:scale-105 active:scale-95"
                                >
                                    Login
                                </Button>
                            </NuxtLink>
                        </NavigationMenuItem>
                        <NavigationMenuItem v-if="!authStore.loggedIn">
                            <NuxtLink to="/register">
                                <Button
                                    variant="default"
                                    size="sm"
                                    class="dark transition-all duration-200 hover:scale-105 hover:shadow-md active:scale-95"
                                >
                                    Register
                                </Button>
                            </NuxtLink>
                        </NavigationMenuItem>
                        <NavigationMenuItem v-if="authStore.loggedIn" class="ml-4">
                            <NuxtLink to="/profile">
                                <Button
                                    variant="link"
                                    class="dark group transition-all duration-200 hover:scale-105 active:scale-95"
                                >
                                    <Icon name="lucide:user" class="mr-2 h-4 w-4 transition-transform duration-200 group-hover:scale-110" aria-hidden="true" />
                                    My Profile
                                </Button>
                            </NuxtLink>
                        </NavigationMenuItem>
                    </NavigationMenuList>
                </NavigationMenu>

                <!-- Mobile Search Toggle -->
                <Button
                    variant="ghost"
                    size="icon"
                    class="md:hidden dark transition-transform duration-150 hover:scale-110 active:scale-95"
                    @click="toggleMobileSearch"
                    :aria-label="showMobileSearch ? 'Close search' : 'Open search'"
                    :aria-expanded="showMobileSearch"
                    aria-controls="mobile-search"
                >
                    <Icon
                        :name="showMobileSearch ? 'lucide:x' : 'lucide:search'"
                        class="h-5 w-5"
                        aria-hidden="true"
                    />
                </Button>
            </div>
        </div>

        <!-- Mobile Search (Expandable) -->
        <Transition
            enter-active-class="transition-all duration-200 ease-out"
            enter-from-class="opacity-0 -translate-y-2"
            enter-to-class="opacity-100 translate-y-0"
            leave-active-class="transition-all duration-150 ease-in"
            leave-from-class="opacity-100 translate-y-0"
            leave-to-class="opacity-0 -translate-y-2"
        >
            <div
                v-if="showMobileSearch"
                id="mobile-search"
                class="container mx-auto mt-4 md:hidden"
            >
                <form
                    role="search"
                    aria-label="Mobile search games"
                    @submit.prevent="handleMobileSearch"
                >
                    <div class="relative">
                        <Input
                            id="mobile-search-input"
                            ref="mobileSearchInput"
                            v-model="searchQuery"
                            type="search"
                            placeholder="Search games..."
                            class="dark w-full pr-10"
                            aria-label="Search for games"
                            :aria-busy="isSearching"
                        />
                        <Button
                            type="submit"
                            variant="ghost"
                            size="icon"
                            class="absolute right-0 top-0 h-full px-3 dark transition-transform duration-150 hover:scale-110 active:scale-95"
                            :disabled="isSearching || !searchQuery.trim()"
                            :aria-label="isSearching ? 'Searching...' : 'Submit search'"
                        >
                            <Icon
                                v-if="!isSearching"
                                name="lucide:search"
                                class="h-4 w-4"
                                aria-hidden="true"
                            />
                            <Icon
                                v-else
                                name="lucide:loader-2"
                                class="h-4 w-4 animate-spin"
                                aria-hidden="true"
                            />
                            <span class="sr-only">{{ isSearching ? 'Searching...' : 'Search' }}</span>
                        </Button>
                    </div>
                </form>
            </div>
        </Transition>
    </nav>
</template>

<script lang="ts" setup>
import { navigationMenuTriggerStyle } from "@/components/ui/navigation-menu";

const authStore = useAuthStore();
const searchQuery = ref("");
const isSearching = ref(false);
const showMobileSearch = ref(false);
const mobileSearchInput = ref<HTMLInputElement | null>(null);

const handleSearch = async () => {
    if (searchQuery.value.trim()) {
        isSearching.value = true;
        try {
            await navigateTo({
                path: "/search",
                query: {
                    s: searchQuery.value,
                },
            });
        } finally {
            // Reset after a brief delay to show the loading state
            setTimeout(() => {
                isSearching.value = false;
            }, 300);
        }
    }
};

const handleMobileSearch = async () => {
    await handleSearch();
    showMobileSearch.value = false;
};

const toggleMobileSearch = () => {
    showMobileSearch.value = !showMobileSearch.value;

    // Focus the mobile search input when opened
    if (showMobileSearch.value) {
        nextTick(() => {
            mobileSearchInput.value?.focus();
        });
    }
};

// Close mobile search on escape key
onMounted(() => {
    const handleEscape = (e: KeyboardEvent) => {
        if (e.key === 'Escape' && showMobileSearch.value) {
            showMobileSearch.value = false;
        }
    };

    document.addEventListener('keydown', handleEscape);

    onUnmounted(() => {
        document.removeEventListener('keydown', handleEscape);
    });
});
</script>

<style scoped>
/* Add any additional styles here if needed */
</style>
