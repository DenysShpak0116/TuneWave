export type CreateTrackRequest = {
    userID: string;
    title: string;
    genre: string;
    artists: string[];
    tags: string[];
    song: File;
    cover: File;
}