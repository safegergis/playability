export default defineEventHandler(async (event) => {
    // Get JWT from cookie
    const jwt = getCookie(event, "jwt");

    if (!jwt) {
        console.log("[my-reports.get] No JWT cookie found");
        throw createError({
            statusCode: 401,
            statusMessage: "Not authenticated - no JWT token found",
        });
    }

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    console.log("[my-reports.get] Fetching user reports");

    const res = await event.$fetch(`${backendUrl}/user/my-reports`, {
        method: "GET",
        headers: {
            Authorization: `Bearer ${jwt}`,
        },
    }).catch((error) => {
        console.log("[my-reports.get] Error fetching user reports:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || "Failed to fetch user reports",
        });
    });

    console.log("[my-reports.get] User reports fetched successfully");
    return res;
});
