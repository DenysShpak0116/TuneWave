export interface CollectiveRankEntry {
    songId: string;
    songName: string;
}

export type CollectiveRank = Record<string, CollectiveRankEntry>;

export type ProfileTable = Record<
    string,
    Record<
        string, 
        string[]
    >
>;

export interface CollectiveDecisionResponse {
    collectiveRank: CollectiveRank;
    profileTable: ProfileTable;
}