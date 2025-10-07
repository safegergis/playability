<template>
    <div
        class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-900 via-green-900 to-slate-900">
        <div class="max-w-2xl mx-auto px-6 text-center">
            <div class="space-y-8">
                <!-- Icon or Logo -->
                <div class="flex justify-center">
                    <div class="w-34 h-34 rounded-full flex items-center justify-center">
                        <img src="/playability_logo.png" alt="Playability logo" class="h-30 w-30 object-contain"
                            width="40" height="40" />
                    </div>
                </div>

                <!-- Main Heading -->
                <div class="space-y-4">
                    <h1 class="text-5xl font-bold text-white">
                        Playability
                    </h1>
                    <h2 class="text-2xl font-semibold text-green-300">
                        Coming Soon
                    </h2>
                </div>

                <!-- Description -->
                <p class="text-lg text-slate-300 max-w-xl mx-auto">
                    We're currently setting up our email verification system to ensure the best experience for our
                    users.
                    User registration and login will be available shortly.
                </p>

                <!-- Status Badge -->
                <div
                    class="inline-flex items-center gap-2 px-6 py-3 bg-green-500/10 border border-green-500/20 rounded-full">
                    <div class="w-2 h-2 bg-green-400 rounded-full animate-pulse"></div>
                    <span class="text-green-300 font-medium">Setting up mail verification</span>
                </div>

                <!-- Additional Info -->
                <div class="pt-8 text-sm text-slate-400">
                    <p>Thank you for your patience as we finalize our authentication system.</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
// Temporary maintenance page for auth routes
</script>
