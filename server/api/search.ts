import type { SearchApiResponse } from "./apiTypes";

const config = useRuntimeConfig();
const igdbClientSecret = config.igdbClientSecret;

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const searchTerm = query.s as string;

  const data = await $fetch<SearchApiResponse>(
    "https://api.igdb.com/v4/search",
    {
      method: "POST",
      headers: {
        "Client-ID": "7bzjkp4ruewaj55ofw7atbugy7p0au",
        Authorization: `Bearer ${igdbClientSecret}`,
      },
      body: `fields id,game,name; search "${searchTerm}"; limit 50;`,
    }
  );

  return data;
});
