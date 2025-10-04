<template>
  <div class="min-h-screen flex items-center justify-center">
    <Card class="w-1/2 min-w-96 shadow-md rounded-md p-6 dark">
      <CardHeader>
        <CardTitle class="text-2xl font-semibold text-center">
          My Profile
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div v-if="authStore.loggedIn && authStore.user" class="space-y-4">
          <div class="flex flex-col space-y-2">
            <span class="text-sm font-medium text-gray-400">Username</span>
            <span class="text-lg">{{ authStore.user.username }}</span>
          </div>
          <div class="flex flex-col space-y-2">
            <span class="text-sm font-medium text-gray-400">User ID</span>
            <span class="text-lg">{{ authStore.user.id }}</span>
          </div>
          <div class="flex flex-col space-y-2">
            <span class="text-sm font-medium text-gray-400">
              Number of Reports
            </span>
            <span class="text-lg">{{ authStore.user.num_of_reports }}</span>
          </div>
          <div class="mt-6">
            <Button variant="destructive" class="w-full" @click="handleLogout">
              Logout
            </Button>
          </div>
        </div>
        <div v-else class="text-center">
          <p class="text-gray-400 mb-4">You are not logged in.</p>
          <NuxtLink to="/login">
            <Button>Go to Login</Button>
          </NuxtLink>
        </div>
      </CardContent>
    </Card>
  </div>
</template>

<script lang="ts" setup>
const authStore = useAuthStore();

const handleLogout = () => {
  authStore.logUserOut();
  navigateTo("/");
};
</script>
