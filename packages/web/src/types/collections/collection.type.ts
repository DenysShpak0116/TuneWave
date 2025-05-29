import { ISong } from "types/song/song.type";

export interface ICollection {
    id: string;
    title: string;
    description: string;
    coverUrl: string;
    createdAt: string;
    user: {
        id: string;
        username: string;
        profilePictureUrl: string;
        profileInfo: string;
    }
    collectionSongs: ISong[]
}