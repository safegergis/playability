declare global {
  interface Cover {
    id: number;
    image_id: string;
  }

  interface Game {
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
    color_blind?: boolean;
    closed_captions?: boolean;
    full_controller_support?: boolean;
    controller_remapping?: boolean;
  }

  interface SearchResult {
    id: number;
    name: string;
  }
}

export {};
