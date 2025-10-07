export default defineEventHandler(async (event) => {
    const query = getQuery(event);
    const searchQuery = query.search;

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    try {
        const response = await $fetch(`${backendUrl}/search`,
            {
                method: 'GET',
                query: { search: searchQuery },
            });
        return response;
    } catch (error: any) {
        console.error('[search.get] Error fetching search results:', error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || 'Failed to fetch search results',
        });
    }
});
