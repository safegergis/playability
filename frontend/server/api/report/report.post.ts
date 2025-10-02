export default defineEventHandler(async (event) => {
  const body = await readBody(event);
  const jwt = getCookie(event, "jwt");

  const res = await $fetch("http://localhost:8080/user/report", {
    method: "POST",
    headers: {
      Authorization: `Bearer ${jwt}`,
    },
    body: body,
  });
  console.log(res);
});
