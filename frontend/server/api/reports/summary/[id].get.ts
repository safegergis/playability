export default defineEventHandler(async (event) => {
    const gameId = getRouterParam(event, 'id');
    if (!gameId) {
        throw createError({
            statusCode: 400,
            statusMessage: 'Game ID is required',
        });
    }

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    try {
        const response = await $fetch(`${backendUrl}/reports/summary/${gameId}`, {
            method: 'GET',
        });
        return response;
    } catch (error: any) {
        if (error.statusCode === 404) {
            return null; // Return null if not found
        }
        console.error(`[reports/summary/${gameId}.get] Error fetching summary:`, error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch summary',
        });
    }
});
