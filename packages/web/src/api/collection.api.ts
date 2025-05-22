import { GetCollectionParams } from "types/collections/getCollectionParams";
import { $api } from "./base.api";

export const createCollection = async (formData: FormData) => {
    try {
        const { data } = await $api.post("/collections", formData, {
            headers: {
                "Content-Type": "multipart/form-data"
            },
        })
        return data;
    }
    catch (err) {
        console.log(err);
    }
}

export const getByUserId = async () => {
    const { data } = await $api.get("/collections/users-collections")
    return data;
}

export const getCollections = async (params: Partial<GetCollectionParams> = {}) => {
    const { data } = await $api.get("/collections/", { params })
    return data
}

export const getTopCollections = async () => {
    const { data } = await $api.get("/collections/")
    return data
}

export const getCollectionById = async (id: string) => {
    const { data } = await $api.get(`/collections/${id}`)
    return data;
}