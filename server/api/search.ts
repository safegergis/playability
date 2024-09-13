import apicalypse from "apicalypse";
const config = useRuntimeConfig();
const igdbClientSecret = config.igdbClientSecret;

const requestOptions = {
  method: "post",
  baseURL: "https://api.igdb.com/v4/games",
  headers: {
    "Client-ID": "7bzjkp4ruewaj55ofw7atbugy7p0au",
    Authorization: `Bearer ${igdbClientSecret}`,
  },
};

export default defineEventHandler(async (event) => {
  const query = getQuery(event);

  const searchTerm = query.s as string;
});
