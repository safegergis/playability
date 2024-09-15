import { useAsyncData } from "nuxt/app";

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const gameID = query.gameID as string;

  const data = await (<string>(`http://localhost:8080/games/${gameID}`,
  {
    method: "GET",
  }));

  /*
  const platformDefinitions = [
    { id: 48, name: "Playstation 4", icon: "mdi:sony-playstation" },
    { id: 167, name: "Playstation 5", icon: "mdi:sony-playstation" },
    { id: 49, name: "Xbox One", icon: "mdi:microsoft-xbox" },
    { id: 169, name: "Xbox Series X|S", icon: "mdi:microsoft-xbox" },
    { id: 6, name: "PC", icon: "mdi:steam" },
    { id: 130, name: "Nintendo Switch", icon: "mdi:nintendo-switch" },
  ];

  console.log("Raw platforms data:", gameData[0].platforms);

  const platforms = gameData[0].platforms
    ?.map((platform) => {
      return platformDefinitions.find(
        (definition) => definition.id === (platform as unknown as number)
      );
    })
    .filter(Boolean);

  console.log("Mapped platforms:", platforms);
*/
  console.log("game data fetched:", data);
});
