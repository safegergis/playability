<template>
  <div class="min-h-screen flex items-center justify-center p-4 dark">
    <!-- Profile Card -->
    <Transition
      mode="out-in"
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="opacity-0 scale-95 -translate-y-4"
      enter-to-class="opacity-100 scale-100 translate-y-0"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="opacity-100 scale-100 translate-y-0"
      leave-to-class="opacity-0 scale-95 translate-y-4"
    >
      <Card
        v-if="!isLoading"
        class="w-full max-w-2xl shadow-lg rounded-xl p-6 dark transition-all duration-300 hover:shadow-xl"
      >
        <CardHeader>
          <div class="flex items-center justify-center gap-3">
            <div
              class="flex h-12 w-12 items-center justify-center rounded-full bg-primary/10 transition-all duration-200 hover:scale-110"
              aria-hidden="true"
            >
              <Icon name="lucide:user" class="h-6 w-6 text-primary" />
            </div>
            <CardTitle class="text-3xl font-bold">
              My Profile
            </CardTitle>
          </div>
        </CardHeader>

        <CardContent>
          <!-- Logged In State -->
          <div v-if="authStore.loggedIn && authStore.user" class="space-y-6">
            <!-- Profile Information -->
            <div
              class="space-y-4 rounded-lg border border-border/50 bg-muted/20 p-6 transition-all duration-200 hover:border-border hover:bg-muted/30"
              role="region"
              aria-label="Profile information"
            >
              <!-- Username with Donator Badge -->
              <div class="group flex flex-col space-y-2 transition-all duration-200">
                <label class="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors duration-200">
                  Username
                </label>
                <div class="flex items-center gap-3">
                  <Icon name="lucide:at-sign" class="h-4 w-4 text-muted-foreground" aria-hidden="true" />
                  <span class="text-lg font-semibold">{{ authStore.user.username }}</span>

                  <!-- TODO: Implement donator badge logic -->
                  <!-- Donator Badge (shown when user is a donator) -->
                  <Transition
                    enter-active-class="transition-all duration-200 ease-out"
                    enter-from-class="opacity-0 scale-75"
                    enter-to-class="opacity-100 scale-100"
                  >
                    <div
                      v-if="isDonator"
                      class="inline-flex items-center gap-1.5 rounded-full bg-gradient-to-r from-amber-500 to-orange-500 px-3 py-1 text-xs font-semibold text-white shadow-md transition-all duration-200 hover:scale-105 hover:shadow-lg"
                      role="status"
                      aria-label="Donator supporter badge"
                    >
                      <Icon name="lucide:heart" class="h-3 w-3" aria-hidden="true" />
                      <span>Supporter</span>
                    </div>
                  </Transition>
                </div>
              </div>

              <!-- Number of Reports -->
              <div class="group flex flex-col space-y-2 transition-all duration-200">
                <label class="text-sm font-medium text-muted-foreground group-hover:text-foreground transition-colors duration-200">
                  Accessibility Reports Submitted
                </label>
                <div class="flex items-center gap-2">
                  <Icon name="lucide:file-text" class="h-4 w-4 text-muted-foreground" aria-hidden="true" />
                  <span class="text-2xl font-bold text-primary">{{ authStore.user.num_of_reports }}</span>
                  <span class="text-sm text-muted-foreground">
                    {{ authStore.user.num_of_reports === 1 ? 'report' : 'reports' }}
                  </span>
                </div>
              </div>
            </div>

            <!-- TODO: Implement payment code redemption -->
            <!-- Donator Section (Payment Code Redemption) -->
            <div
              v-if="!isDonator"
              class="rounded-lg border border-primary/20 bg-gradient-to-br from-primary/5 to-primary/10 p-6 transition-all duration-200 hover:border-primary/30 hover:from-primary/10 hover:to-primary/15"
              role="region"
              aria-label="Supporter benefits"
            >
              <div class="mb-4 flex items-start gap-3">
                <div class="flex h-10 w-10 items-center justify-center rounded-full bg-primary/20 transition-all duration-200 hover:scale-110" aria-hidden="true">
                  <Icon name="lucide:heart" class="h-5 w-5 text-primary" />
                </div>
                <div class="flex-1">
                  <h3 class="text-lg font-semibold">Become a Supporter</h3>
                  <p class="text-sm text-muted-foreground">
                    Support accessibility and get a special badge
                  </p>
                </div>
              </div>

              <!-- Payment Code Form -->
              <form
                @submit.prevent="redeemPaymentCode"
                class="space-y-3"
                role="form"
                aria-label="Redeem supporter code"
              >
                <div class="space-y-2">
                  <label
                    for="payment-code"
                    class="text-sm font-medium"
                  >
                    Enter Payment Code
                  </label>
                  <div class="relative">
                    <Input
                      id="payment-code"
                      v-model="paymentCode"
                      type="text"
                      placeholder="XXXX-XXXX-XXXX"
                      class="dark pr-10 font-mono transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                      :disabled="isRedeeming"
                      :aria-busy="isRedeeming"
                      :aria-invalid="!!paymentError"
                      :aria-describedby="paymentError ? 'payment-error' : undefined"
                      autocomplete="off"
                    />
                    <Icon
                      v-if="paymentSuccess"
                      name="lucide:check-circle"
                      class="absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-green-500 transition-all duration-200"
                      aria-hidden="true"
                    />
                  </div>

                  <!-- Error Message -->
                  <Transition
                    enter-active-class="transition-all duration-200 ease-out"
                    enter-from-class="opacity-0 -translate-y-1"
                    enter-to-class="opacity-100 translate-y-0"
                    leave-active-class="transition-all duration-150 ease-in"
                    leave-from-class="opacity-100 translate-y-0"
                    leave-to-class="opacity-0 -translate-y-1"
                  >
                    <p
                      v-if="paymentError"
                      id="payment-error"
                      class="text-sm text-destructive"
                      role="alert"
                    >
                      {{ paymentError }}
                    </p>
                  </Transition>
                </div>

                <Button
                  type="submit"
                  class="w-full transition-all duration-200 hover:scale-[1.02] hover:shadow-md active:scale-[0.98] focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
                  :disabled="!paymentCode.trim() || isRedeeming"
                  :aria-busy="isRedeeming"
                >
                  <Icon
                    v-if="!isRedeeming"
                    name="lucide:gift"
                    class="mr-2 h-4 w-4"
                    aria-hidden="true"
                  />
                  <Icon
                    v-else
                    name="lucide:loader-2"
                    class="mr-2 h-4 w-4 animate-spin"
                    aria-hidden="true"
                  />
                  {{ isRedeeming ? 'Redeeming...' : 'Redeem Code' }}
                </Button>
              </form>
            </div>

            <!-- Actions -->
            <div class="flex flex-col gap-3 pt-2" role="group" aria-label="Profile actions">
              <!-- View Reports Button -->
              <NuxtLink
                to="/my-reports"
                class="group"
              >
                <Button
                  variant="outline"
                  class="w-full dark transition-all duration-200 hover:scale-[1.02] hover:shadow-md active:scale-[0.98] focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
                >
                  <Icon name="lucide:list" class="mr-2 h-4 w-4 transition-transform duration-200 group-hover:scale-110" aria-hidden="true" />
                  View My Reports
                </Button>
              </NuxtLink>

              <!-- Logout Button -->
              <Button
                variant="destructive"
                class="w-full transition-all duration-200 hover:scale-[1.02] hover:shadow-md active:scale-[0.98] focus-visible:ring-2 focus-visible:ring-destructive focus-visible:ring-offset-2"
                :disabled="isLoggingOut"
                @click="confirmLogout"
                :aria-busy="isLoggingOut"
              >
                <Icon
                  v-if="!isLoggingOut"
                  name="lucide:log-out"
                  class="mr-2 h-4 w-4"
                  aria-hidden="true"
                />
                <Icon
                  v-else
                  name="lucide:loader-2"
                  class="mr-2 h-4 w-4 animate-spin"
                  aria-hidden="true"
                />
                {{ isLoggingOut ? 'Logging out...' : 'Logout' }}
              </Button>
            </div>
          </div>

          <!-- Not Logged In State -->
          <div
            v-else
            class="space-y-6 text-center py-8"
            role="status"
            aria-live="polite"
          >
            <div class="flex flex-col items-center gap-4">
              <div
                class="flex h-16 w-16 items-center justify-center rounded-full bg-muted transition-all duration-200"
                aria-hidden="true"
              >
                <Icon name="lucide:user-x" class="h-8 w-8 text-muted-foreground" />
              </div>
              <div class="space-y-2">
                <p class="text-lg font-medium text-muted-foreground">
                  You are not logged in
                </p>
                <p class="text-sm text-muted-foreground">
                  Please log in to view your profile and accessibility reports
                </p>
              </div>
            </div>
            <NuxtLink to="/login" class="inline-block">
              <Button
                size="lg"
                class="transition-all duration-200 hover:scale-105 hover:shadow-md active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
              >
                <Icon name="lucide:log-in" class="mr-2 h-4 w-4" aria-hidden="true" />
                Go to Login
              </Button>
            </NuxtLink>
          </div>
        </CardContent>
      </Card>

      <!-- Loading Skeleton -->
      <Card
        v-else
        class="w-full max-w-2xl shadow-lg rounded-xl p-6 dark"
        aria-busy="true"
        aria-label="Loading profile"
      >
        <CardHeader>
          <div class="flex items-center justify-center gap-3">
            <div class="h-12 w-12 rounded-full bg-muted animate-pulse" />
            <div class="h-8 w-32 bg-muted animate-pulse rounded" />
          </div>
        </CardHeader>
        <CardContent class="space-y-6">
          <div class="space-y-4 rounded-lg border border-border/50 bg-muted/20 p-6">
            <div class="space-y-2">
              <div class="h-4 w-20 bg-muted animate-pulse rounded" />
              <div class="h-6 w-40 bg-muted animate-pulse rounded" />
            </div>
            <div class="space-y-2">
              <div class="h-4 w-16 bg-muted animate-pulse rounded" />
              <div class="h-6 w-32 bg-muted animate-pulse rounded" />
            </div>
            <div class="space-y-2">
              <div class="h-4 w-48 bg-muted animate-pulse rounded" />
              <div class="h-6 w-24 bg-muted animate-pulse rounded" />
            </div>
          </div>
          <div class="space-y-3">
            <div class="h-10 w-full bg-muted animate-pulse rounded" />
            <div class="h-10 w-full bg-muted animate-pulse rounded" />
          </div>
        </CardContent>
      </Card>
    </Transition>

    <!-- Logout Confirmation Dialog -->
    <Dialog v-model:open="showLogoutDialog">
      <DialogContent class="dark">
        <DialogHeader>
          <DialogTitle class="flex items-center gap-2">
            <Icon name="lucide:log-out" class="h-5 w-5 text-destructive" aria-hidden="true" />
            Confirm Logout
          </DialogTitle>
          <DialogDescription>
            Are you sure you want to log out? You will need to log in again to access your profile and submit reports.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter class="gap-2 sm:gap-0">
          <Button
            variant="outline"
            class="transition-all duration-200 hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
            @click="showLogoutDialog = false"
          >
            Cancel
          </Button>
          <Button
            variant="destructive"
            class="transition-all duration-200 hover:scale-105 hover:shadow-md active:scale-95 focus-visible:ring-2 focus-visible:ring-destructive focus-visible:ring-offset-2"
            @click="handleLogout"
          >
            <Icon name="lucide:log-out" class="mr-2 h-4 w-4" aria-hidden="true" />
            Logout
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script lang="ts" setup>
// Define page metadata
definePageMeta({
  title: 'My Profile',
});

// Use SEO meta tags
useHead({
  title: 'My Profile - Playability',
  meta: [
    {
      name: 'description',
      content: 'View your Playability profile and accessibility report contributions',
    },
  ],
});

const authStore = useAuthStore();

// State management
const isLoading = ref(true);
const isLoggingOut = ref(false);
const showLogoutDialog = ref(false);

// TODO: Implement donator status from backend/database
// Donator/Supporter state
const isDonator = ref(false); // TODO: Get from user data or API
const paymentCode = ref('');
const isRedeeming = ref(false);
const paymentError = ref('');
const paymentSuccess = ref(false);

// Simulate initial load (can be replaced with actual data fetching)
onMounted(() => {
  // Small delay to show loading skeleton
  setTimeout(() => {
    isLoading.value = false;
  }, 300);
});

// TODO: Implement payment code redemption logic
// Redeem payment code
const redeemPaymentCode = async () => {
  // Clear previous states
  paymentError.value = '';
  paymentSuccess.value = false;

  // Validate code format (basic validation)
  if (!paymentCode.value.trim()) {
    paymentError.value = 'Please enter a payment code';
    return;
  }

  isRedeeming.value = true;

  try {
    // TODO: Replace with actual API call
    // Example: await $fetch('/api/redeem-payment-code', { method: 'POST', body: { code: paymentCode.value } })

    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500));

    // TODO: Handle actual API response
    // For now, simulate an error for demonstration
    throw new Error('TODO: Implement payment code redemption API');

    // On success:
    // paymentSuccess.value = true;
    // isDonator.value = true;
    // paymentCode.value = '';
    // Show success toast/notification
  } catch (error: any) {
    console.error('Payment code redemption error:', error);
    paymentError.value = error.message || 'Invalid payment code. Please try again.';
  } finally {
    isRedeeming.value = false;
  }
};

// Show confirmation dialog before logout
const confirmLogout = () => {
  showLogoutDialog.value = true;
};

// Handle logout with loading state
const handleLogout = async () => {
  isLoggingOut.value = true;

  try {
    // Add a small delay for UX (shows loading state)
    await new Promise(resolve => setTimeout(resolve, 500));

    // Perform logout
    authStore.logUserOut();

    // Navigate to home
    await navigateTo('/');
  } catch (error) {
    console.error('Logout error:', error);
  } finally {
    isLoggingOut.value = false;
    showLogoutDialog.value = false;
  }
};

// Keyboard shortcuts
onMounted(() => {
  const handleKeyPress = (e: KeyboardEvent) => {
    // Close dialog with Escape key
    if (e.key === 'Escape' && showLogoutDialog.value) {
      showLogoutDialog.value = false;
    }
  };

  document.addEventListener('keydown', handleKeyPress);

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeyPress);
  });
});
</script>
