export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch(`${backendUrl}/user/request-password-reset`, {
        method: "POST",
        body: body,
    }).catch((error) => {
        console.log("[request-password-reset.post] Error calling backend:", error);
        // Don't throw error for security - always return success
        // to prevent email enumeration
    });

    console.log("[request-password-reset.post] Password reset request processed");

    return { success: true, message: "If the email exists, a password reset code has been sent" };
});
