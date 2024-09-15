export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const searchTerm = query.s as string;

  const data = await $fetch<string>(
    `http://localhost:8080/search/${searchTerm}`,
    {
      method: "GET",
    }
  );

  return JSON.parse(data);
});
