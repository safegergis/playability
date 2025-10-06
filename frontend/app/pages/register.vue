<template>
    <main
        class="min-h-screen dark flex items-center justify-center bg-gradient-to-b from-background via-background to-neutral-900 py-12 px-4">
        <div class="w-full max-w-md animate-scale-in">
            <!-- Logo/Icon Section -->
            <div class="text-center mb-8 animate-fade-in">
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-4 transition-transform duration-300 hover:scale-110">
                    <Icon name="lucide:user-plus" class="w-10 h-10 text-4xl text-primary" aria-hidden="true" />
                </div>
                <h1 class="text-3xl font-bold">Create Account</h1>
                <p class="text-muted-foreground mt-2">Join Playability to start tracking accessible games</p>
            </div>

            <Card class="dark shadow-2xl transition-all duration-300 hover:shadow-3xl animate-slide-up"
                style="animation-delay: 100ms">
                <CardHeader class="space-y-1 pb-4">
                    <CardTitle class="text-2xl font-bold text-center">Register</CardTitle>
                    <CardDescription class="text-center">
                        Create your account to get started
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <form @submit.prevent="onSubmit" class="space-y-4">
                        <!-- Username Field -->
                        <FormField v-slot="{ componentField }" name="username">
                            <FormItem>
                                <FormLabel for="username" class="text-sm font-medium">
                                    Username
                                    <span class="text-destructive ml-1" aria-label="required">*</span>
                                </FormLabel>
                                <FormDescription class="text-xs text-muted-foreground">
                                    This is your public display name
                                </FormDescription>
                                <FormControl>
                                    <div class="relative">
                                        <Icon name="lucide:user"
                                            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none"
                                            aria-hidden="true" />
                                        <Input v-bind="componentField" id="username" type="text" autocomplete="username"
                                            placeholder="johndoe"
                                            class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                                            :aria-invalid="!!(form.errors.value.username || uniqueUsername)"
                                            :aria-describedby="(form.errors.value.username || uniqueUsername) ? 'username-error' : undefined"
                                            @input="uniqueUsername = false" />
                                    </div>
                                </FormControl>
                                <div v-if="form.errors.value.username"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    id="username-error" role="alert">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ form.errors.value.username }}</span>
                                </div>
                                <div v-if="uniqueUsername"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    role="alert" aria-live="assertive">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>Username is already in use</span>
                                </div>
                            </FormItem>
                        </FormField>

                        <!-- Email Field -->
                        <FormField v-slot="{ componentField }" name="email">
                            <FormItem>
                                <FormLabel for="email" class="text-sm font-medium">
                                    Email Address
                                    <span class="text-destructive ml-1" aria-label="required">*</span>
                                </FormLabel>
                                <FormControl>
                                    <div class="relative">
                                        <Icon name="lucide:mail"
                                            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none"
                                            aria-hidden="true" />
                                        <Input v-bind="componentField" id="email" type="email" autocomplete="email"
                                            placeholder="you@example.com"
                                            class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                                            :aria-invalid="!!(form.errors.value.email || uniqueEmail)"
                                            :aria-describedby="(form.errors.value.email || uniqueEmail) ? 'email-error' : undefined"
                                            @input="uniqueEmail = false" />
                                    </div>
                                </FormControl>
                                <div v-if="form.errors.value.email"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    id="email-error" role="alert">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ form.errors.value.email }}</span>
                                </div>
                                <div v-if="uniqueEmail"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    role="alert" aria-live="assertive">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>Email is already in use</span>
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
                                        <Icon name="lucide:lock"
                                            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none"
                                            aria-hidden="true" />
                                        <Input v-bind="componentField" id="password" type="password"
                                            autocomplete="new-password" placeholder="Minimum 6 characters"
                                            class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                                            :aria-invalid="!!form.errors.value.password"
                                            :aria-describedby="form.errors.value.password ? 'password-error' : undefined" />
                                    </div>
                                </FormControl>
                                <div v-if="form.errors.value.password"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    id="password-error" role="alert">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ form.errors.value.password }}</span>
                                </div>
                            </FormItem>
                        </FormField>

                        <!-- Confirm Password Field -->
                        <FormField v-slot="{ componentField }" name="confirmPassword">
                            <FormItem>
                                <FormLabel for="confirmPassword" class="text-sm font-medium">
                                    Confirm Password
                                    <span class="text-destructive ml-1" aria-label="required">*</span>
                                </FormLabel>
                                <FormControl>
                                    <div class="relative">
                                        <Icon name="lucide:lock-keyhole"
                                            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none"
                                            aria-hidden="true" />
                                        <Input v-bind="componentField" id="confirmPassword" type="password"
                                            autocomplete="new-password" placeholder="Re-enter your password"
                                            class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                                            :aria-invalid="!!form.errors.value.confirmPassword"
                                            :aria-describedby="form.errors.value.confirmPassword ? 'confirm-password-error' : undefined" />
                                    </div>
                                </FormControl>
                                <div v-if="form.errors.value.confirmPassword"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    id="confirm-password-error" role="alert">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ form.errors.value.confirmPassword }}</span>
                                </div>
                            </FormItem>
                        </FormField>

                        <!-- Submit Button -->
                        <Button type="submit" :disabled="isSubmitting" :aria-busy="isSubmitting"
                            class="w-full transition-all duration-200 hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2">
                            <span v-if="!isSubmitting" class="flex items-center justify-center gap-2">
                                <Icon name="lucide:user-plus" class="w-5 h-5" aria-hidden="true" />
                                Create Account
                            </span>
                            <span v-else class="flex items-center justify-center gap-2">
                                <Icon name="lucide:loader-2" class="w-5 h-5 animate-spin" aria-hidden="true" />
                                Creating account...
                            </span>
                        </Button>
                    </form>

                    <!-- Divider -->
                    <div class="relative my-6">
                        <div class="absolute inset-0 flex items-center">
                            <span class="w-full border-t border-border"></span>
                        </div>
                        <div class="relative flex justify-center text-xs uppercase">
                            <span class="bg-card px-2 text-muted-foreground">Already have an account?</span>
                        </div>
                    </div>

                    <!-- Login Link -->
                    <div class="text-center">
                        <NuxtLink to="/login"
                            class="inline-flex items-center gap-2 text-sm text-primary hover:text-primary/80 underline-offset-4 transition-all duration-200 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-sm px-2 py-1">
                            <Icon name="lucide:log-in" class="w-4 h-4" aria-hidden="true" />
                            Sign in to existing account
                        </NuxtLink>
                    </div>
                </CardContent>
            </Card>

            <!-- Additional Info -->
            <p class="text-center text-sm text-muted-foreground mt-6 animate-fade-in" style="animation-delay: 200ms">
                By creating an account, you agree to our
                <NuxtLink to="/terms"
                    class="text-primary hover:text-primary/80 underline-offset-4 hover:underline transition-colors">
                    Terms of Service</NuxtLink>
                and
                <NuxtLink to="/privacy"
                    class="text-primary hover:text-primary/80 underline-offset-4 hover:underline transition-colors">
                    Privacy Policy</NuxtLink>
            </p>
        </div>
    </main>
</template>

<script lang="ts" setup>
import { useForm } from "vee-validate";
import * as yup from "yup";

// SEO metadata
useHead({
    title: 'Create Account - Playability',
    meta: [
        {
            name: 'description',
            content: 'Create a Playability account to track accessible games, submit accessibility reports, and help build a more inclusive gaming community.'
        }
    ]
});

const authStore = useAuthStore();

const uniqueEmail = ref(false);
const uniqueUsername = ref(false);
const isSubmitting = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
    username: yup
        .string()
        .required("Username is required")
        .min(3, "Username must be at least 3 characters"),
    email: yup
        .string()
        .required("Email is required")
        .email("Must be a valid email"),
    password: yup
        .string()
        .required("Password is required")
        .min(6, "Password must be at least 6 characters"),
    confirmPassword: yup
        .string()
        .required("Please confirm your password")
        .oneOf([yup.ref("password")], "Passwords must match"),
});

// Initialize useForm with the validation schema
const form = useForm({
    validationSchema: schema,
});

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
    isSubmitting.value = true;
    uniqueEmail.value = false;
    uniqueUsername.value = false;

    const registrationData = {
        username: values.username,
        email: values.email,
        password: values.password,
    };

    try {
        const { data, error } = await useFetch("/api/auth/register", {
            method: "POST",
            body: registrationData,
            immediate: true,
        });

        if (error.value) {
            console.log(error.value.data)
            if (error.value.data.data.includes("email")) {
                uniqueEmail.value = true;
            } else if (error.value.data.data.includes("username")) {
                uniqueUsername.value = true;
            } else {
                //TODO: Error response goes here
            }
        } else {
            // Redirect to email verification page
            await navigateTo(`/verify-email?email=${encodeURIComponent(values.email)}`);
        }
    } finally {
        isSubmitting.value = false;
    }
});
</script>

<style></style>
