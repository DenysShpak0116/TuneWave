import { $api } from "./base.api"


export const getListensByDay = async () => {
    const { data } = await $api.get("/analytics/listens-by-day")
    return data
}

export const getAvarageListens = async () => {
    const { data } = await $api.get("/analytics/avg-listens")
    return data
}

export const getMedianListens = async () => {
    const { data } = await $api.get("/analytics/median-listens")
    return data
}

export const getPopularTracks = async () => {
    const { data } = await $api.get("/analytics/most-popular-track")
    return data
}

export const getPeakAndSilentDay = async () => {
    const { data } = await $api.get("/analytics/peak-silent")
    return data
}