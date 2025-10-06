<template>
    <main
        class="min-h-screen dark flex items-center justify-center bg-gradient-to-b from-background via-background to-neutral-900 py-12 px-4">
        <div class="w-full max-w-md animate-scale-in">
            <!-- Logo/Icon Section -->
            <div class="text-center mb-8 animate-fade-in">
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-4 transition-transform duration-300 hover:scale-110">
                    <Icon name="lucide:key-round" class="w-10 h-10 text-4xl text-primary" aria-hidden="true" />
                </div>
                <h1 class="text-3xl font-bold">Reset Password</h1>
                <p class="text-muted-foreground mt-2">Enter your email to receive a reset code</p>
            </div>

            <Card class="dark shadow-2xl transition-all duration-300 hover:shadow-3xl animate-slide-up"
                style="animation-delay: 100ms">
                <CardHeader class="space-y-1 pb-4">
                    <CardTitle class="text-2xl font-bold text-center">Forgot Password</CardTitle>
                    <CardDescription class="text-center">
                        We'll send you a verification code to reset your password
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <!-- Success Message -->
                    <div v-if="requestSuccess"
                        class="flex items-center gap-2 p-3 mb-4 rounded-lg bg-green-500/10 border border-green-500/20 text-green-500 animate-fade-in"
                        role="alert" aria-live="polite">
                        <Icon name="lucide:mail-check" class="w-5 h-5 flex-shrink-0" aria-hidden="true" />
                        <p class="text-sm font-medium">Reset code sent! Check your email.</p>
                    </div>

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
                                        <Icon name="lucide:mail"
                                            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none"
                                            aria-hidden="true" />
                                        <Input v-bind="componentField" id="email" type="email" autocomplete="email"
                                            placeholder="you@example.com"
                                            class="dark pl-10 transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                                            :aria-invalid="!!form.errors.value.email"
                                            :aria-describedby="form.errors.value.email ? 'email-error' : undefined" />
                                    </div>
                                </FormControl>
                                <div v-if="form.errors.value.email"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    id="email-error" role="alert">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ form.errors.value.email }}</span>
                                </div>
                            </FormItem>
                        </FormField>

                        <!-- Submit Button -->
                        <Button type="submit" :disabled="isSubmitting" :aria-busy="isSubmitting"
                            class="w-full transition-all duration-200 hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2">
                            <span v-if="!isSubmitting" class="flex items-center justify-center gap-2">
                                <Icon name="lucide:send" class="w-5 h-5" aria-hidden="true" />
                                Send Reset Code
                            </span>
                            <span v-else class="flex items-center justify-center gap-2">
                                <Icon name="lucide:loader-2" class="w-5 h-5 animate-spin" aria-hidden="true" />
                                Sending...
                            </span>
                        </Button>
                    </form>

                    <!-- Divider -->
                    <div class="relative my-6">
                        <div class="absolute inset-0 flex items-center">
                            <span class="w-full border-t border-border"></span>
                        </div>
                        <div class="relative flex justify-center text-xs uppercase">
                            <span class="bg-card px-2 text-muted-foreground">Remember your password?</span>
                        </div>
                    </div>

                    <!-- Login Link -->
                    <div class="text-center">
                        <NuxtLink to="/login"
                            class="inline-flex items-center gap-2 text-sm text-primary hover:text-primary/80 underline-offset-4 transition-all duration-200 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 rounded-sm px-2 py-1">
                            <Icon name="lucide:log-in" class="w-4 h-4" aria-hidden="true" />
                            Back to login
                        </NuxtLink>
                    </div>
                </CardContent>
            </Card>

            <!-- Additional Info -->
            <p class="text-center text-sm text-muted-foreground mt-6 animate-fade-in" style="animation-delay: 200ms">
                Need help?
                <NuxtLink to="/support"
                    class="text-primary hover:text-primary/80 underline-offset-4 hover:underline transition-colors">
                    Contact support</NuxtLink>
            </p>
        </div>
    </main>
</template>

<script lang="ts" setup>
import { useForm } from "vee-validate";
import * as yup from "yup";

// SEO metadata
useHead({
    title: 'Forgot Password - Playability',
    meta: [
        {
            name: 'description',
            content: 'Reset your Playability account password.'
        }
    ]
});

const isSubmitting = ref(false);
const requestSuccess = ref(false);

// Define the validation schema using Yup
const schema = yup.object({
    email: yup
        .string()
        .required("Email is required")
        .email("Must be a valid email"),
});

// Initialize useForm with the validation schema
const form = useForm({
    validationSchema: schema,
});

// Handle form submission
const onSubmit = form.handleSubmit(async (values) => {
    isSubmitting.value = true;
    requestSuccess.value = false;

    try {
        const { data, error } = await useFetch("/api/auth/request-password-reset", {
            method: "POST",
            body: {
                email: values.email,
            },
            immediate: true,
        });

        // Always show success for security (don't reveal if email exists)
        requestSuccess.value = true;

        // Redirect to reset password page after a delay
        setTimeout(() => {
            navigateTo(`/reset-password?email=${encodeURIComponent(values.email)}`);
        }, 2000);
    } finally {
        isSubmitting.value = false;
    }
});
</script>

<style></style>
