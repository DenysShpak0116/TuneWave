import { IAuthor } from "./author.type";
import { ISongTags } from "./songTags.type";
import { IUser } from "../user/user.type";

export interface ISong {
    id: string;
    createdAt: string;
    duration: string;
    genre: string;
    songUrl: string;
    coverUrl: string
    listenings: number;
    likes: number;
    dislikes: number;
    user: IUser
    authors: IAuthor[]
    songTags: ISongTags[]
}


