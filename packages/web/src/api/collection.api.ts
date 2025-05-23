import { GetCollectionParams } from "types/collections/getCollectionParams";
import { $api } from "./base.api";
import { CollectionSongsParams } from "types/collections/collectionSongsParams.type";

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

export const updateCollection = async (collectionId: string, formData: FormData) => {
    try {
        const { data } = await $api.put(`/collections/${collectionId}`, formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        })
        return data
    }
    catch (err) {
        console.log(err)
    }
};

export const getCollectionSongs = async (collectionId: string, params: Partial<CollectionSongsParams> = {}) => {
    const { data } = await $api.get(`collections/${collectionId}/songs`, { params });
    return data;
}