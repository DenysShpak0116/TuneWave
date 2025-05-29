import { IVectorType } from "types/vectors/vector.type";
import { $api } from "./base.api";

export const addVector = async (collectionId: string, songId: string, vectors: IVectorType[]) => {
    const { data } = await $api.post(`/collections/${collectionId}/${songId}/vectors`, { vectors })
    return data
}

export const getVectorBySongId = async (collectionId: string, songId: string) => {
    const { data } = await $api.get(`/collections/${collectionId}/${songId}/vectors`)
    return data
}

export const updateVector = async (collectionId: string, songId: string, vectors: IVectorType[]) => {
    const { data } = await $api.put(`collections/${collectionId}/${songId}/vectors`, { vectors })
    return data
}

export const hasAllVectors = async (collectionId: string) => {
    const { data } = await $api.get(`/collections/${collectionId}/has-all-vectors`);
    return data.hasAllVectors as boolean;
}

export const deleteVector = async (collectionId: string, songId: string) => {
    const { data } = await $api.delete(`/collections/${collectionId}/${songId}/vectors`)
    return data
}