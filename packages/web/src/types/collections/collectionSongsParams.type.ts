export interface CollectionSongsParams {
    search: string;
    sortBy: "added_at" | "title";
    order: "asc" | "desc";
    page: number;
    limit: number;
}