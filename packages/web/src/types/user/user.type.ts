import { FollowType } from "./follow.type"

export interface IUser {
    id: string
    username: string
    role: string
    profileInfo: string
    email: string
    profilePictureUrl: string
    follows: FollowType[];
    followers: FollowType[];
}