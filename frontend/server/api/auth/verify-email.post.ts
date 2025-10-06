export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch(`${backendUrl}/user/verify-email`, {
        method: "POST",
        body: body,
    }).catch((error) => {
        console.log("[verify-email.post] Error calling backend:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.data?.message || error.message || "Email verification failed",
            data: error.data,
        });
    });

    console.log("[verify-email.post] Email verified successfully");

    return { success: true, message: "Email verified successfully" };
});
