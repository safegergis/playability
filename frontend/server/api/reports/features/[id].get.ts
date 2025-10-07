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
        const response = await $fetch(`${backendUrl}/reports/features/${gameId}`, {
            method: 'GET',
        });
        return response;
    } catch (error: any) {
        console.error(`[reports/features/${gameId}.get] Error fetching feature stats:`, error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch feature stats',
        });
    }
});
