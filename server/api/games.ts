import type { GameApiResponse, CoverApiResponse } from "./apiTypes";

const config = useRuntimeConfig();
const igdbClientSecret = config.igdbClientSecret;

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const gameID = query.gameID as string;

  const data = await $fetch<GameApiResponse>(`https://api.igdb.com/v4/games`, {
    method: "POST",
    headers: {
      "Client-ID": "7bzjkp4ruewaj55ofw7atbugy7p0au",
      Authorization: `Bearer ${igdbClientSecret}`,
    },
    body: `fields name,cover,summary,platforms,involved_companies; where id = ${gameID};`,
  });

  const coverData = await $fetch<CoverApiResponse>(
    `https://api.igdb.com/v4/covers`,
    {
      method: "POST",
      headers: {
        "Client-ID": "7bzjkp4ruewaj55ofw7atbugy7p0au",
        Authorization: `Bearer ${igdbClientSecret}`,
      },
      body: `fields image_id; where id = ${data[0].cover};`,
    }
  );

  const response = {
    ...data[0],
    cover: `https://images.igdb.com/igdb/image/upload/t_cover_big/${coverData[0].image_id}.jpg`,
  };

  return response;
});
