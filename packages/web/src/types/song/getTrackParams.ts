export interface GetTracksParams {
    search?: string;
    sortBy?: 'created_at' | 'title' | 'artist' | 'genre' | 'listenings';
    order?: 'asc' | 'desc';
    page?: number;
    limit?: number;
}
