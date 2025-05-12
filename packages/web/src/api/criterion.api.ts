import { $api } from "./base.api"

export const getAllCriterion = async () => {
    const { data } = await $api.get("/criterions/")
    return data
}

export const createCriterion = async (name: string) => {
    const { data } = await $api.post("/criterions/", { name })
    return data
}

export const deleteCriterion = async (id: string) => {
    const { data } = await $api.delete(`/criterions/${id}`)
    return data
  }