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
        const response = await $fetch(`${backendUrl}/reports/score/${gameId}`, {
            method: 'GET',
        });
        return response;
    } catch (error: any) {
        if (error.statusCode === 404) {
            return null; // Return null if not found
        }
        console.error(`[reports/score/${gameId}.get] Error fetching score:`, error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch score',
        });
    }
});
