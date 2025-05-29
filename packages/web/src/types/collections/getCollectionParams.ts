export interface GetCollectionParams {
    limit?: number;
    page?: number;
    sort?: "tile" | "created_at";
    order?: "asc"
}