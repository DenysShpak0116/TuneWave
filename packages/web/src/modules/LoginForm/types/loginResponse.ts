import { IUser } from "types/user/user.type"

export interface LoginResponse {
    accessToken: string,
    refreshToken: string
    user: IUser
}