<template>
  <main class="min-h-screen dark flex items-center justify-center bg-gradient-to-b from-background via-background to-neutral-900 py-12 px-4">
    <div class="w-full max-w-md animate-scale-in">
      <!-- Logo/Icon Section -->
      <div class="text-center mb-8 animate-fade-in">
        <div class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-4 transition-transform duration-300 hover:scale-110">
          <Icon name="lucide:gamepad-2" class="w-10 h-10 text-5xl text-primary" aria-hidden="true" />
        </div>
        <h1 class="text-3xl font-bold">Welcome Back</h1>
        <p class="text-muted-foreground mt-2">Sign in to your Playability account</p>
      </div>

      <Card class="dark shadow-2xl transition-all duration-300 hover:shadow-3xl animate-slide-up" style="animation-delay: 100ms">
        <CardHeader class="space-y-1 pb-4">
          <CardTitle class="text-2xl font-bold text-center">Login</CardTitle>
          <CardDescription class="text-center">
            Enter your credentials to access your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="onSubmit" class="space-y-4">
            <!-- Email Field -->
            <FormField v-slot="{ componentField }" name="email">
              <FormItem>
                <FormLabel for="email" class="text-sm font-medium">
                  Email Address
                  <span class="text-destructive ml-1" aria-label="required">*</span>
                </FormLabel>
                <FormControl>
                  <div class="relative">
                    <Icon name="lucide:mail" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none" aria-hidden="true" />
                    <Input
                      v-bind="componentField"
                      id="email"
                      type="email"
                      autocomplete="email"
                      placeholder="you@example.com"
                      class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                      :aria-invalid="!!form.errors.value.email"
                      :aria-describedby="form.errors.value.email ? 'email-error' : undefined"
                      @input="InvalidLogin = false"
                    />
                  </div>
                </FormControl>
                <div v-if="form.errors.value.email" class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake" id="email-error" role="alert">
                  <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                  <span>{{ form.errors.value.email }}</span>
                </div>
              </FormItem>
            </FormField>

            <!-- Password Field -->
            <FormField v-slot="{ componentField }" name="password">
              <FormItem>
                <FormLabel for="password" class="text-sm font-medium">
                  Password
                  <span class="text-destructive ml-1" aria-label="required">*</span>
                </FormLabel>
                <FormControl>
                  <div class="relative">
                    <Icon name="lucide:lock" class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none" aria-hidden="true" />
                    <Input
                      v-bind="componentField"
                      id="password"
                      type="password"
                      autocomplete="current-password"
                      placeholder="Enter your password"
                      class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                      :aria-invalid="!!form.errors.value.password"
                      :aria-describedby="form.errors.value.password ? 'password-error' : undefined"
                      @input="InvalidLogin = false"
                    />
                  </div>
                </FormControl>
                <div v-if="form.errors.value.password" class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake" id="password-error" role="alert">
                  <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                  <span>{{ form.errors.value.password }}</span>
                </div>
              </FormItem>
            </FormField>

            <!-- Login Error Message -->
            <div
              v-if="InvalidLogin"
              class="flex items-center gap-2 p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive animate-shake"
              role="alert"
              aria-live="assertive"
            >
              <Icon name="lucide:alert-triangle" class="w-5 h-5 flex-shrink-0" aria-hidden="true" />
              <p class="text-sm font-medium">Invalid email or password. Please try again.</p>
            </div>

            <!-- Submit Button -->
            <Button
              type="submit"
              :disabled="isSubmitting"
              :aria-busy="isSubmitting"
              class="w-full transition-all duration-200 hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2"
            >
              <span v-if="!isSubmitting" class="flex items-center justify-center gap-2">
                <Icon name="lucide:log-in" class="w-5 h-5" aria-hidden="true" />
                Sign In
              </span>
              <span v-else class="flex items-center justify-center gap-2">
                <Icon name="lucide:loader-2" class="w-5 h-5 animate-spin" aria-hidden="true" />
                Signing in...
              </span>
            </Button>
          </form>

          <!-- Divider -->
          <div class="relative my-6">
            <div class="absolute inset-0 flex items-center">
              <span class="w-full border-t border-border"></span>
            </div>
            <div class="relative flex justify-center text-xs uppercase">
              <span class="bg-card px-2 text-muted-foreground">New to Playability?</span>
            </div>
          </div>

          <!-- Register Link -->
          <div class="text-center">
            <NuxtLink
              to="/register"
              class="inline-flex items-center gap-2 text-sm text-primary hover:text-primary/80 underline-offset-4 transition-all duration-200 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-sm px-2 py-1"
            >
              <Icon name="lucide:user-plus" class="w-4 h-4" aria-hidden="true" />
              Create an account
            </NuxtLink>
          </div>
        </CardContent>
      </Card>

      <!-- Additional Info -->
      <p class="text-center text-sm text-muted-foreground mt-6 animate-fade-in" style="animation-delay: 200ms">
        By signing in, you agree to our
        <NuxtLink to="/terms" class="text-primary hover:text-primary/80 underline-offset-4 hover:underline transition-colors">Terms of Service</NuxtLink>
        and
        <NuxtLink to="/privacy" class="text-primary hover:text-primary/80 underline-offset-4 hover:underline transition-colors">Privacy Policy</NuxtLink>
      </p>
    </div>
  </main>
</template>

<script lang="ts" setup>
import { useForm } from "vee-validate";
import * as yup from "yup";

// SEO metadata
useHead({
  title: 'Sign In - Playability',
  meta: [
    {
      name: 'description',
      content: 'Sign in to your Playability account to access game accessibility information and submit reports.'
    }
  ]
});

const authStore = useAuthStore();

const InvalidLogin = ref(false);
const isSubmitting = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
  email: yup
    .string()
    .required("Email is required")
    .email("Must be a valid email"),
  password: yup
    .string()
    .required("Password is required")
    .min(6, "Password must be at least 6 characters"),
});

// Initialize useForm with the validation schema
const form = useForm({
  validationSchema: schema,
});

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
  isSubmitting.value = true;
  InvalidLogin.value = false;

  try {
    const { data, error } = await useFetch("/api/auth/login", {
      method: "POST",
      body: values,
      immediate: true,
    });
    if (error.value?.statusCode === 401) {
      InvalidLogin.value = true;
    } else if (error.value) {
      console.error("Login error:", error.value);
      InvalidLogin.value = true;
    } else {
      authStore.logUserIn(data.value!.user);
      await navigateTo("/");
    }
  } finally {
    isSubmitting.value = false;
  }
});
</script>

<style></style>
