export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    // Debug: log all cookies and headers
    const allCookies = parseCookies(event);
    console.log("[report.post] All cookies:", allCookies);
    console.log("[report.post] Request headers:", getHeaders(event));

    const jwt = getCookie(event, "jwt");

    if (!jwt) {
        console.log("[report.post] No JWT cookie found");
        throw createError({
            statusCode: 401,
            statusMessage: "Not authenticated - no JWT token found",
        });
    }

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    console.log("report.post jwt: " + jwt);
    const res = await event.$fetch(`${backendUrl}/user/report`, {
        method: "POST",
        body: body,
        async onRequest({ request, options }) {
            // Log the request body before sending
            console.log("Request URL:", request);
            console.log("Request Method:", options.method);
            console.log("Request Headers:", options.headers);
            console.log("Request Body:", options.body); // Log the body
        },
    }).catch((error) => {
        console.log("[report.post] Error submitting report:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || "Failed to submit report",
        });
    });

    console.log("[report.post] Report submitted successfully");
    return res;
});
