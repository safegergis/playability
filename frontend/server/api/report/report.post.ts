export default defineEventHandler(async (event) => {
    const body = await readBody(event);
    const jwt = getCookie(event, "jwt");

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch(`${backendUrl}/user/report`, {
        method: "POST",
        headers: {
            Authorization: `Bearer ${jwt}`,
        },
        body: body,
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
