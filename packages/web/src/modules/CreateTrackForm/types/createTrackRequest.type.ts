export type CreateTrackRequest = {
    userId: string;
    title: string;
    genre: string;
    artists: string[];
    tags: string[];
    song: File;
    cover: File;
}