export interface Cover {
  id: number;
  image_id: string;
}

export interface Game {
  id: number;
  name: string;
  cover?:
    | {
        id: number;
        url: string;
      }
    | string;
  summary?: string;
  platforms?: number[];
  involved_companies?: {
    id: number;
    company: {
      id: number;
      name: string;
    };
  }[];
  cover_art?: string;
}

export interface SearchResult {
  id: number;
  name: string;
}
