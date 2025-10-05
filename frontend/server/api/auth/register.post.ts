export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch(`${backendUrl}/user/register`, {
        method: "POST",
        body: body,
    }).catch((error) => {
        console.log("[register.post] Error calling backend:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || "Registration failed",
            data: error.data,
        });
    });

    console.log("[register.post] Registration successful, setting cookie");

    return { success: true };
});
