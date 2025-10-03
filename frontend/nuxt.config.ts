// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: "2024-04-03",
    devtools: { enabled: false },
    modules: [
        "@nuxt/eslint",
        "@nuxtjs/tailwindcss",
        "shadcn-nuxt",
        "@nuxt/icon",
        "@nuxt/image",
        "@vee-validate/nuxt",
        "@pinia/nuxt",
    ],

    runtimeConfig: {
        // Private keys (server-side only)
        apiUrl: process.env.NUXT_API_URL || 'http://backend:8080',

        // Public keys (exposed to client)
        public: {
            apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8080',
        },
    },

    tailwindcss: {
        cssPath: "~/assets/css/tailwind.css",
        configPath: "~/tailwind.config.js",
    },
    shadcn: {
        /**
         * Prefix for all the imported component
         */
        prefix: "",
        /**
         * Directory that the component lives in.
         * @default "./components/ui"
         */
        componentDir: "./components/ui",
    },
    vite: {
        optimizeDeps: {
            exclude: ["vee-validate"],
        },
    },
});
