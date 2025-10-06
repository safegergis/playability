export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch(`${backendUrl}/user/resend-verification`, {
        method: "POST",
        body: body,
    }).catch((error) => {
        console.log("[resend-verification.post] Error calling backend:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.data?.message || error.message || "Failed to resend verification email",
            data: error.data,
        });
    });

    console.log("[resend-verification.post] Verification email resent successfully");

    return { success: true, message: "Verification email sent" };
});
