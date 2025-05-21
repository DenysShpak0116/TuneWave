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

export const getChatPreviews = async() => {
    return (await $api.get("/chats")).data.chats
}

export const updateUser = async (id: string, profileInfo: string, username: string) => {
    const { data } = await $api.put(`/users/${id}`, { username: username, profileInfo: profileInfo })
    return data
}
