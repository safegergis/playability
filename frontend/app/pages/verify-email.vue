<template>
    <main
        class="min-h-screen dark flex items-center justify-center bg-gradient-to-b from-background via-background to-neutral-900 py-12 px-4">
        <div class="w-full max-w-md animate-scale-in">
            <!-- Logo/Icon Section -->
            <div class="text-center mb-8 animate-fade-in">
                <div
                    class="inline-flex items-center justify-center w-20 h-20 rounded-full bg-primary/10 mb-4 transition-transform duration-300 hover:scale-110">
                    <Icon name="lucide:mail-check" class="w-10 h-10 text-4xl text-primary" aria-hidden="true" />
                </div>
                <h1 class="text-3xl font-bold">Verify Your Email</h1>
                <p class="text-muted-foreground mt-2">Enter the verification code sent to your email</p>
            </div>

            <Card class="dark shadow-2xl transition-all duration-300 hover:shadow-3xl animate-slide-up"
                style="animation-delay: 100ms">
                <CardHeader class="space-y-1 pb-4">
                    <CardTitle class="text-2xl font-bold text-center">Email Verification</CardTitle>
                    <CardDescription class="text-center">
                        We've sent a 6-digit code to <strong>{{ email }}</strong>
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <!-- Success Message -->
                    <div v-if="verificationSuccess"
                        class="flex items-center gap-2 p-3 mb-4 rounded-lg bg-green-500/10 border border-green-500/20 text-green-500 animate-fade-in"
                        role="alert" aria-live="polite">
                        <Icon name="lucide:check-circle" class="w-5 h-5 flex-shrink-0" aria-hidden="true" />
                        <p class="text-sm font-medium">Email verified! Redirecting to login...</p>
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
                                <div v-if="invalidCode"
                                    class="inline-flex items-center gap-1 text-sm text-destructive mt-1 animate-shake"
                                    role="alert" aria-live="assertive">
                                    <Icon name="lucide:alert-circle" class="w-3 h-3 flex-shrink-0" aria-hidden="true" />
                                    <span>{{ errorMessage }}</span>
                                </div>
                            </FormItem>
                        </FormField>

                        <!-- Submit Button -->
                        <Button type="submit" :disabled="isSubmitting" :aria-busy="isSubmitting"
                            class="w-full transition-all duration-200 hover:scale-105 active:scale-95 focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2">
                            <span v-if="!isSubmitting" class="flex items-center justify-center gap-2">
                                <Icon name="lucide:check" class="w-5 h-5" aria-hidden="true" />
                                Verify Email
                            </span>
                            <span v-else class="flex items-center justify-center gap-2">
                                <Icon name="lucide:loader-2" class="w-5 h-5 animate-spin" aria-hidden="true" />
                                Verifying...
                            </span>
                        </Button>
                    </form>

                    <!-- Divider -->
                    <div class="relative my-6">
                        <div class="absolute inset-0 flex items-center">
                            <span class="w-full border-t border-border"></span>
                        </div>
                        <div class="relative flex justify-center text-xs uppercase">
                            <span class="bg-card px-2 text-muted-foreground">Didn't receive the code?</span>
                        </div>
                    </div>

                    <!-- Resend Code Button -->
                    <Button type="button" variant="outline" @click="resendCode" :disabled="isResending || cooldown > 0"
                        class="w-full transition-all duration-200">
                        <span v-if="!isResending" class="flex items-center justify-center gap-2">
                            <Icon name="lucide:refresh-cw" class="w-4 h-4" aria-hidden="true" />
                            {{ cooldown > 0 ? `Resend code in ${cooldown}s` : 'Resend code' }}
                        </span>
                        <span v-else class="flex items-center justify-center gap-2">
                            <Icon name="lucide:loader-2" class="w-4 h-4 animate-spin" aria-hidden="true" />
                            Sending...
                        </span>
                    </Button>

                    <!-- Success message for resend -->
                    <div v-if="resendSuccess"
                        class="flex items-center gap-2 p-3 mt-4 rounded-lg bg-green-500/10 border border-green-500/20 text-green-500 animate-fade-in"
                        role="alert" aria-live="polite">
                        <Icon name="lucide:mail" class="w-4 h-4 flex-shrink-0" aria-hidden="true" />
                        <p class="text-sm">Verification code sent!</p>
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
    title: 'Verify Email - Playability',
    meta: [
        {
            name: 'description',
            content: 'Verify your email address to complete your Playability account setup.'
        }
    ]
});

const route = useRoute();
const email = ref((route.query.email as string) || '');

const invalidCode = ref(false);
const errorMessage = ref('Invalid or expired verification code');
const isSubmitting = ref(false);
const isResending = ref(false);
const verificationSuccess = ref(false);
const resendSuccess = ref(false);
const cooldown = ref(0);

// Redirect if no email provided
onMounted(() => {
    if (!email.value) {
        navigateTo('/register');
    }
});

// Cooldown timer
const startCooldown = () => {
    cooldown.value = 60;
    const interval = setInterval(() => {
        cooldown.value--;
        if (cooldown.value <= 0) {
            clearInterval(interval);
        }
    }, 1000);
};

// Define the validation schema using Yup
const schema = yup.object({
    code: yup
        .string()
        .required("Verification code is required")
        .matches(/^\d{6}$/, "Code must be 6 digits"),
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
        const { data, error } = await useFetch("/api/auth/verify-email", {
            method: "POST",
            body: {
                email: email.value,
                code: values.code,
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
            verificationSuccess.value = true;
            setTimeout(() => {
                navigateTo("/login");
            }, 2000);
        }
    } finally {
        isSubmitting.value = false;
    }
});

// Handle resend code
const resendCode = async () => {
    isResending.value = true;
    resendSuccess.value = false;

    try {
        const { error } = await useFetch("/api/auth/resend-verification", {
            method: "POST",
            body: {
                email: email.value,
            },
            immediate: true,
        });

        if (error.value) {
            invalidCode.value = true;
            errorMessage.value = 'Failed to resend code. Please try again.';
        } else {
            resendSuccess.value = true;
            startCooldown();
            setTimeout(() => {
                resendSuccess.value = false;
            }, 5000);
        }
    } finally {
        isResending.value = false;
    }
};
</script>

<style></style>
