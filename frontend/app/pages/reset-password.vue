<template>
    <main
        class="min-h-screen dark flex items-center justify-center bg-gradient-to-b from-background via-background to-neutral-900 py-12 px-4">
        <div class="w-full max-w-md animate-scale-in">
            <!-- Logo/Icon Section -->
            <div class="text-center mb-8 animate-fade-in">
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-4 transition-transform duration-300 hover:scale-110">
                    <Icon name="lucide:lock-keyhole" class="w-10 h-10 text-4xl text-primary" aria-hidden="true" />
                </div>
                <h1 class="text-3xl font-bold">Reset Password</h1>
                <p class="text-muted-foreground mt-2">Enter the code and your new password</p>
            </div>

            <Card class="dark shadow-2xl transition-all duration-300 hover:shadow-3xl animate-slide-up"
                style="animation-delay: 100ms">
                <CardHeader class="space-y-1 pb-4">
                    <CardTitle class="text-2xl font-bold text-center">Create New Password</CardTitle>
                    <CardDescription class="text-center">
                        Enter the code sent to <strong>{{ email }}</strong>
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <!-- Success Message -->
                    <div v-if="resetSuccess"
                        class="flex items-center gap-2 p-3 mb-4 rounded-lg bg-green-500/10 border border-green-500/20 text-green-500 animate-fade-in"
                        role="alert" aria-live="polite">
                        <Icon name="lucide:check-circle" class="w-5 h-5 flex-shrink-0" aria-hidden="true" />
                        <p class="text-sm font-medium">Password reset! Redirecting to login...</p>
                    </div>

                    <form @submit.prevent="onSubmit" class="space-y-4">
                        <!-- Verification Code Field -->
                        <FormField v-slot="{ componentField }" name="code">
                            <FormItem>
                                <FormLabel for="code" class="text-sm font-medium">
                                    Verification Code
                                    <span class="text-destructive ml-1" aria-label="required">*</span>
                                </FormLabel>
                                <FormControl>
                                    <div class="relative">
                                        <Icon name="lucide:key"
                                            class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-muted-foreground pointer-events-none"
                                            aria-hidden="true" />
                                        <Input v-bind="componentField" id="code" type="text" autocomplete="off"
                                            placeholder="Enter 6-digit code"
                                            class="dark pl-10 text-center text-lg tracking-widest transition-all duration-200 focus-visible:ring-2 focus-visible:ring-primary"
                                            :aria-invalid="!!(form.errors.value.code || invalidCode)"
                                            :aria-describedby="(form.errors.value.code || invalidCode) ? 'code-error' : undefined"
                                            @input="invalidCode = false" maxlength="6" />
                                    </div>
                                </FormControl>
                                <div v-if="form.errors.value.code"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    id="code-error" role="alert">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ form.errors.value.code }}</span>
                                </div>
                            </FormItem>
                        </FormField>

                        <!-- New Password Field -->
                        <FormField v-slot="{ componentField }" name="password">
                            <FormItem>
                                <FormLabel for="password" class="text-sm font-medium">
                                    New Password
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
                                    Confirm New Password
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

                        <!-- Error Message -->
                        <div v-if="invalidCode"
                            class="flex items-center gap-2 p-3 rounded-lg bg-destructive/10 border border-destructive/20 text-destructive animate-shake"
                            role="alert" aria-live="assertive">
                            <Icon name="lucide:alert-triangle" class="w-5 h-5 flex-shrink-0" aria-hidden="true" />
                            <p class="text-sm font-medium">{{ errorMessage }}</p>
                        </div>

                        <!-- Submit Button -->
                        <Button type="submit" :disabled="isSubmitting" :aria-busy="isSubmitting"
                            class="w-full transition-all duration-200 hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2">
                            <span v-if="!isSubmitting" class="flex items-center justify-center gap-2">
                                <Icon name="lucide:check" class="w-5 h-5" aria-hidden="true" />
                                Reset Password
                            </span>
                            <span v-else class="flex items-center justify-center gap-2">
                                <Icon name="lucide:loader-2" class="w-5 h-5 animate-spin" aria-hidden="true" />
                                Resetting...
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
    title: 'Reset Password - Playability',
    meta: [
        {
            name: 'description',
            content: 'Create a new password for your Playability account.'
        }
    ]
});

const route = useRoute();
const email = ref((route.query.email as string) || '');

const invalidCode = ref(false);
const errorMessage = ref('Invalid or expired verification code');
const isSubmitting = ref(false);
const resetSuccess = ref(false);

// Redirect if no email provided
onMounted(() => {
    if (!email.value) {
        navigateTo('/forgot-password');
    }
});

// Define the validation schema using Yup
const schema = yup.object({
    code: yup
        .string()
        .required("Verification code is required")
        .matches(/^\d{6}$/, "Code must be 6 digits"),
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
    invalidCode.value = false;

    try {
        const { data, error } = await useFetch("/api/auth/reset-password", {
            method: "POST",
            body: {
                email: email.value,
                code: values.code,
                new_password: values.password,
            },
            immediate: true,
        });

        if (error.value) {
            invalidCode.value = true;
            if (error.value.statusCode === 400) {
                errorMessage.value = error.value.data?.message || 'Invalid or expired verification code';
            } else {
                errorMessage.value = 'An error occurred. Please try again.';
            }
        } else {
            resetSuccess.value = true;
            setTimeout(() => {
                navigateTo("/login");
            }, 2000);
        }
    } finally {
        isSubmitting.value = false;
    }
});
</script>

<style></style>
