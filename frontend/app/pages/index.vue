<template>
    <main class="container mx-auto px-4 py-8 lg:py-12 space-y-16">
        <!-- Hero Section -->
        <section aria-labelledby="hero-heading" class="relative overflow-hidden">
            <div
                class="relative rounded-2xl bg-gradient-to-br from-primary/10 via-primary/5 to-transparent border border-primary/20 p-8 lg:p-8 -mb-10">

                <!-- Content -->
                <div class="relative z-5 dark text-center max-w-3xl mx-auto ">
                    <div
                        class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-primary/10 border border-primary/20 text-sm font-medium text-primary mb-2">
                        <Icon name="lucide:trending-up" class="w-4 h-4" aria-hidden="true" />
                        <span>Trending Now</span>
                    </div>
                    <h1 id="hero-heading"
                        class="text-4xl lg:text-5xl xl:text-6xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-transparent">
                        Discover What's Popular
                    </h1>
                    <p class="text-lg lg:text-xl text-muted-foreground max-w-2xl mx-auto leading-relaxed">
                        Explore the most played games right now and see how accessible they are for everyone
                    </p>
                </div>
            </div>
        </section>

        <!-- Featured Games Section -->
        <section aria-labelledby="featured-games-heading" class="animate-slide-up">
            <h2 id="featured-games-heading"
                class="text-3xl lg:text-4xl font-bold mb-6 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Games Being Played
            </h2>
            <GameCarousel :games="FeaturedGames" :status="status" />
        </section>

        <!-- Get Started Section -->
        <section aria-labelledby="get-started-heading" class="animate-slide-up" style="animation-delay: 150ms">
            <h2 id="get-started-heading"
                class="text-3xl lg:text-4xl font-bold mb-8 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
                Get Started
            </h2>
            <div class="grid md:grid-cols-3 gap-6 lg:gap-8">
                <!-- Search Games Card -->
                <Card
                    class="group dark transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 cursor-pointer"
                    tabindex="0" role="article" aria-label="Search games for accessibility information"
                    @click="navigateTo('/search')" @keydown.enter="navigateTo('/search')"
                    @keydown.space.prevent="navigateTo('/search')">
                    <CardHeader>
                        <CardTitle class="flex items-center gap-3 text-xl">
                            <Icon name="lucide:search"
                                class="w-8 h-8 text-primary transition-transform duration-300 group-hover:scale-110 group-hover:rotate-12"
                                aria-hidden="true" />
                            <span>Search Games</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p class="text-muted-foreground leading-relaxed">
                            Find accessibility information for your favorite games.
                        </p>
                    </CardContent>
                </Card>

                <!-- Submit Reports Card -->
                <Card
                    class="group dark transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 cursor-pointer"
                    tabindex="0" role="article" aria-label="Submit accessibility reports" @click="navigateTo('/submit')"
                    @keydown.enter="navigateTo('/submit')" @keydown.space.prevent="navigateTo('/submit')">
                    <CardHeader>
                        <CardTitle class="flex items-center gap-3 text-xl">
                            <Icon name="lucide:clipboard-pen"
                                class="w-8 h-8 text-primary transition-transform duration-300 group-hover:scale-110 group-hover:-rotate-12"
                                aria-hidden="true" />
                            <span>Submit Reports</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p class="text-muted-foreground leading-relaxed">
                            Share your experiences and help others by submitting accessibility reports.
                        </p>
                    </CardContent>
                </Card>

                <!-- Join Community Card -->
                <Card
                    class="group dark transition-all duration-300 hover:-translate-y-2 hover:shadow-2xl focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2 cursor-pointer"
                    tabindex="0" role="article" aria-label="Join the Playability community"
                    @click="navigateTo('/community')" @keydown.enter="navigateTo('/community')"
                    @keydown.space.prevent="navigateTo('/community')">
                    <CardHeader>
                        <CardTitle class="flex items-center gap-3 text-xl">
                            <Icon name="lucide:users"
                                class="w-8 h-8 text-primary transition-all duration-300 group-hover:scale-110"
                                aria-hidden="true" />
                            <span>Join the Community</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p class="text-muted-foreground leading-relaxed">
                            Connect with other gamers and share insights on game accessibility.
                        </p>
                    </CardContent>
                </Card>
            </div>
        </section>

        <!-- Support & Mission Section -->
        <section aria-labelledby="support-mission-heading" class="grid md:grid-cols-2 gap-8 animate-slide-up"
            style="animation-delay: 300ms">
            <h2 id="support-mission-heading" class="sr-only">Support and Mission</h2>

            <!-- Donate Card -->
            <div
                class="group dark rounded-xl p-6 lg:p-8 bg-gradient-to-br from-card to-card/80 border border-border transition-all duration-300 hover:shadow-xl hover:-translate-y-1 focus-within:ring-2 focus-within:ring-primary focus-within:ring-offset-2">
                <h3 class="text-2xl lg:text-3xl font-bold mb-4 flex items-center gap-2">
                    <Icon name="lucide:heart"
                        class="w-7 h-7 text-red-500 transition-transform duration-300 group-hover:scale-110 group-hover:animate-pulse"
                        aria-hidden="true" />
                    <span>Support Us</span>
                </h3>
                <p class="text-muted-foreground leading-relaxed mb-6">
                    Playability is built and maintained by a single person. If you find it useful, please consider
                    donating to help keep it running.
                </p>
                <a href="https://ko-fi.com/playability" target="_blank" rel="noopener noreferrer"
                    class="inline-block focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-lg"
                    aria-label="Donate on Ko-fi (opens in new tab)">
                    <Button variant="outline"
                        class="dark transition-all duration-200 hover:scale-105 hover:bg-primary hover:text-primary-foreground hover:border-primary active:scale-95">
                        <Icon name="lucide:coffee" class="w-4 h-4 mr-2" aria-hidden="true" />
                        Donate on Ko-fi
                    </Button>
                </a>
            </div>

            <!-- Mission Card -->
            <div
                class="group dark rounded-xl p-6 lg:p-8 bg-gradient-to-br from-card to-card/80 border border-border transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
                <h3 class="text-2xl lg:text-3xl font-bold mb-4 flex items-center gap-2">
                    <Icon name="lucide:target"
                        class="w-7 h-7 text-primary transition-transform duration-300 group-hover:scale-110 group-hover:rotate-12"
                        aria-hidden="true" />
                    <span>Our Mission</span>
                </h3>
                <p class="text-muted-foreground leading-relaxed">
                    Our mission is to create a more inclusive gaming community by highlighting accessibility features in
                    games and providing a platform for users to share their experiences.
                </p>
            </div>
        </section>
    </main>
</template>

<script setup lang="ts">
const config = useRuntimeConfig();

// Set page metadata for SEO and accessibility
useHead({
    title: 'Playability - Game Accessibility Information',
    meta: [
        {
            name: 'description',
            content: 'Discover accessibility features in video games. Search games, read reports, and share your experiences to help create a more inclusive gaming community.'
        }
    ]
});

const { data: FeaturedGames, status } = await useFetch<FeaturedGame[]>(
    '/featured',
    {
        key: 'featured-games',
        baseURL: import.meta.server ? config.apiUrl : config.public.apiUrl,
        default: () => []
    }
);
</script>

<style></style>
