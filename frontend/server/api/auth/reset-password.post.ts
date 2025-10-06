export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch(`${backendUrl}/user/reset-password`, {
        method: "POST",
        body: body,
    }).catch((error) => {
        console.log("[reset-password.post] Error calling backend:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.data?.message || error.message || "Password reset failed",
            data: error.data,
        });
    });

    console.log("[reset-password.post] Password reset successfully");

    return { success: true, message: "Password reset successfully" };
});
