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

  interface FeatureStat {
    name: string;
    true_percentage: number;
    limited_percentage: number;
    false_percentage: number;
    consensus: string;
    secondary_consensus: string;
  }
  interface Report {
    id: string;
    game_id: string;
    user_id: string;
    score: number;
    report: string;
    username?: string;
  }
  interface User {
    id: number;
    username: string;
    email: string;
    hash: string;
    num_of_reports: number;
  }
}

export {};
