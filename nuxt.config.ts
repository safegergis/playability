// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2024-04-03",
  devtools: { enabled: true },
  modules: [
    "@nuxt/eslint",
    "@nuxtjs/tailwindcss",
    "shadcn-nuxt",
    "@nuxt/icon",
    "@nuxt/image",
  ],
  runtimeConfig: {
    igdbClientSecret: process.env.IGDB_CLIENT_SECRET,
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
});