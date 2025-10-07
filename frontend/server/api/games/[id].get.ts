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
        const response = await $fetch(`${backendUrl}/games`, {
            method: 'GET',
            query: { id: gameId },
        });
        return response;
    } catch (error: any) {
        console.error(`[games/${gameId}.get] Error fetching game data:`, error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch game data',
        });
    }
});