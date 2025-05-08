export interface GetTracksParams {
    search?: string;
    sortBy?: 'created_at' | 'title' | 'artist' | 'genre';
    order?: 'asc' | 'desc';
    page?: number;
    limit?: number;
}
