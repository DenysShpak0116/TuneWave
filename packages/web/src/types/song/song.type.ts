import { IAuthor } from "./author.type";
import { ISongTags } from "./songTags.type";
import { IUser } from "../user/user.type";
import { IComment } from "types/comments/comment.type";

export interface ISong {
    id: string;
    createdAt: string;
    title: string;
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
    comments: IComment[]
}


