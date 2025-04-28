import { $api } from "./base.api"

export const getUserById = async (id: string) => {
    const { data } = await $api.get(`/users/${id}`)
    return data
}

export const updateUser = async (id: string, profileInfo: string, username: string) => {
    const { data } = await $api.post("/users", { id, profileInfo, username })
    return data
} 