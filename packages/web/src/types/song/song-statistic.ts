import { IAuthor } from "./author.type";

export interface SongStatistic {
    authors: IAuthor[]
    id: string
    duration: string
    title: string
    genre: string
    songUrl: string
    coverUrl: string
    listenings: number
    likes: number,
    dislikes: number
}

export interface ITopTrackStatistic {
    song: SongStatistic
    count: number
}