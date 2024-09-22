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
    color_blind?: string;
    closed_captions?: string;
    full_controller_support?: string;
    controller_remapping?: string;
    accessibility_score?: number;
  }

  interface SearchResult {
    id: number;
    name: string;
  }
}

export {};
