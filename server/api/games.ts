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
  const platformDefinitions = [
    { id: 48, name: "Playstation 4", icon: "mdi:sony-playstation" },
    { id: 167, name: "Playstation 5", icon: "mdi:sony-playstation" },
    { id: 49, name: "Xbox One", icon: "mdi:microsoft-xbox" },
    { id: 169, name: "Xbox Series X|S", icon: "mdi:microsoft-xbox" },
    { id: 6, name: "PC", icon: "mdi:steam" },
    { id: 130, name: "Nintendo Switch", icon: "mdi:nintendo-switch" },
  ];

  console.log("Raw platforms data:", data[0].platforms);

  const platforms = data[0].platforms
    ?.map((platform) => {
      return platformDefinitions.find(
        (definition) => definition.id === (platform as unknown as number)
      );
    })
    .filter(Boolean);

  console.log("Mapped platforms:", platforms);

  const response = {
    ...data[0],
    coverArt: `https://images.igdb.com/igdb/image/upload/t_cover_big/${coverData[0].image_id}.jpg`,
    platforms: platforms,
  };

  return response;
});
