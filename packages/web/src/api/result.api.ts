import { IResultType } from "types/results/result.type";
import { $api } from "./base.api";

export const createResult = async (collectionId: string, results: IResultType[]) => {
    const { data } = await $api.post(`/collections/${collectionId}/send-results`, { results })
    return data
}

export const getUserResults = async (collectionId: string) => {
    const { data } = await $api.get(`/collections/${collectionId}/get-user-results`)
    return data;
}