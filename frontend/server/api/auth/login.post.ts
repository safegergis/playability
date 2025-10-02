export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    // Use backend service name in Docker, localhost for development
    const backendUrl = useRuntimeConfig(event).public.apiUrl || "http://backend:8080";

    const res = await $fetch(`${backendUrl}/user/login`, {
        method: "POST",
        body: body,
    }).catch((error) => {
        console.log("[login.post] Error calling backend:", error);
        throw createError({
            statusCode: error.statusCode || 500,
            statusMessage: error.message || "Login failed",
        });
    });

    console.log("[login.post] Login successful, setting cookie");
    setCookie(event, "jwt", res as string, {
        httpOnly: true,
        maxAge: 86400,
        path: "/",
        secure: false, // Set to false for development, true for production
        expires: new Date(Date.now() + 86400 * 1000),
    });

    return res;
});
