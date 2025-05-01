import { $api } from "./base.api"

export const getUserById = async (id: string) => {
    const { data } = await $api.get(`/users/${id}`)
    return data
}

export const updateUser = async (id: string, profileInfo: string, username: string) => {
    const { data } = await $api.put(`/users/${id}`, { username: username, profileInfo: profileInfo })
    return data
}
