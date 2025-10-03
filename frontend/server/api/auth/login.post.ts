export default defineEventHandler(async (event) => {
    const body = await readBody(event);

    const config = useRuntimeConfig(event);
    const backendUrl = config.apiUrl;

    const res = await $fetch<LoginResponse>(`${backendUrl}/user/login`, {
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
    console.log("[login.post] JWT token:", res.token);
    setCookie(event, "jwt", res.token, {
        httpOnly: true,
        maxAge: 86400,
        path: "/",
        secure: false, // Set to false for development, true for production
        sameSite: "lax",
    });
    console.log("[login.post] Cookie set, returning response");

    return { success: true, user: res.user };
});
