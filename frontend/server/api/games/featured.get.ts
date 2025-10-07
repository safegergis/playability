export default defineEventHandler(async (event) => {
    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    try {
        const response = await $fetch(`${backendUrl}/featured`, {
            method: 'GET',
        });
        return response;
    } catch (error: any) {
        console.error('[featured.get] Error fetching featured games:', error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch featured games',
        });
    }
});
