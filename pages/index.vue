<template>
  <div class="container mx-auto px-4 py-8">
    <section class="mb-12 p-4">
      <h2 class="text-2xl font-semibold mb-4">Games Being Played</h2>
      <GameCarousel :games="FeaturedGames" :status="status" />
    </section>

    <section class="mb-12 p-4">
      <h2 class="text-2xl font-semibold mb-4">Get Started</h2>
      <div class="grid md:grid-cols-3 gap-6">
        <Card class="dark">
          <CardHeader>
            <CardTitle class="flex items-center">
              <Icon name="lucide:search" class="w-8 h-8 mr-2" />
              Search Games
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p>Find accessibility information for your favorite games.</p>
          </CardContent>
        </Card>
        <Card class="dark">
          <CardHeader>
            <CardTitle class="flex items-center">
              <Icon name="lucide:clipboard-pen" class="w-8 h-8 mr-2" />
              Submit Reports
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Share your experiences and help others by submitting accessibility
              reports.
            </p>
          </CardContent>
        </Card>
        <Card class="dark">
          <CardHeader>
            <CardTitle class="flex items-center">
              <Icon name="lucide:users" class="w-8 h-8 mr-2" />
              Join the Community
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p>
              Connect with other gamers and share insights on game
              accessibility.
            </p>
          </CardContent>
        </Card>
      </div>
    </section>

    <section class="grid md:grid-cols-2 gap-8 p-4">
      <div class="rounded-lg p-4 bg-gray-800">
        <h2 class="text-2xl font-semibold mb-4">Donate</h2>
        <p class="text-gray-300">
          Playability is built and maintained by a single person. If you find it
          useful, please consider donating to help keep it running.
        </p>
        <a href="https://ko-fi.com/playability" target="_blank">
          <Button variant="outline" class="dark mt-4">Donate</Button>
        </a>
      </div>
      <div class="rounded-lg p-4 bg-gray-800">
        <h2 class="text-2xl font-semibold mb-4">Our Mission</h2>
        <p class="text-gray-300">
          Our mission is to create a more inclusive gaming community by
          highlighting accessibility features in games and providing a platform
          for users to share their experiences.
        </p>
      </div>
    </section>
    <Toaster />
  </div>
</template>

<script setup lang="ts">
const FeaturedGames = ref<FeaturedGame[]>([]);

const { data, error, status } = await useLazyAsyncData<FeaturedGame[]>(
  "featured",
  () => $fetch("http://localhost:8080/featured")
);
if (error.value) {
  console.error(error.value);
} else if (data.value) {
  console.log(data.value);
  FeaturedGames.value = data.value;
}
</script>

<style></style>
