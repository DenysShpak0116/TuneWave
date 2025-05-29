import { $api } from "./base.api"

export const getUserById = async (id: string) => {
    const { data } = await $api.get(`/users/${id}`)
    return data
}

export const updateUserAvatar = async (formData: FormData) => {
    try {
        const { data } = await $api.put("/users/avatar/", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            }
        })
        return data
    }
    catch (err) {
        console.log(err)
    }
};

export const getChatPreviews = async () => {
    return (await $api.get("/chats")).data.chats
}

export const updateUser = async (id: string, profileInfo: string, username: string) => {
    const { data } = await $api.put(`/users/${id}`, { username: username, profileInfo: profileInfo })
    return data
}

export const getUserCollections = async (userId: string) => {
    const { data } = await $api.get(`/users/${userId}/collections`);
    return data;
}

export const followUser = async (userId: string) => {
    return (await $api.post(`/users/${userId}/follow`))
}

export const isFollowed = async (userId: string) => {
    return (await $api.get(`/users/${userId}/is-followed`)).data.isFollowed
}

export const unfollowUser = async (userId: string) => {
    return (await $api.delete(`/users/${userId}/unfollow`))
}