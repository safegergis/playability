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
        const response = await $fetch(`${backendUrl}/reports/cards/${gameId}`, {
            method: 'GET',
        });
        return response;
    } catch (error: any) {
        if (error.statusCode === 404) {
            return []; // Return empty array if not found
        }
        console.error(`[reports/cards/${gameId}.get] Error fetching report cards:`, error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch report cards',
        });
    }
});
