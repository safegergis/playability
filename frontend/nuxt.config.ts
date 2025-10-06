// https://nuxt.com/docs/api/configuration/nuxt-config
import tailwindcss from "@tailwindcss/vite";
export default defineNuxtConfig({
    compatibilityDate: "2024-04-03",
    devtools: { enabled: false },
    modules: [
        "@nuxt/eslint",
        "shadcn-nuxt",
        "@nuxt/icon",
        "@nuxt/image",
        "@vee-validate/nuxt",
        "@pinia/nuxt",
    ],
    components: [
        {
            path: '~/components',
            pathPrefix: false,
        }
    ],
    css: ['~/assets/css/main.css'],
    runtimeConfig: {
        // Private keys (server-side only)
        apiUrl: process.env.NUXT_API_URL || 'http://backend:8080',

        // Public keys (exposed to client)
        public: {
            apiUrl: process.env.NUXT_PUBLIC_API_URL || 'http://localhost:8080',
        },
    },

    pinia: {
        storesDirs: ['./stores/**']
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
        componentDir: "./app/components/ui",
    },
    vite: {
        plugins: [
            tailwindcss()
        ],
    },
});
