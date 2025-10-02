export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  const res = await $fetch("http://localhost:8080/user/login", {
    method: "POST",
    body: body,
  }).catch((error) => {
    console.log(error);
    throw createError({
      statusCode: error.statusCode,
      statusMessage: error.message,
    });
  });
  console.log(res);
  setCookie(event, "jwt", res as string, {
    httpOnly: true,
    maxAge: 86400,
    path: "/",
    secure: true,
    expires: new Date(Date.now() + 86400 * 1000),
  });

  return res;
});
