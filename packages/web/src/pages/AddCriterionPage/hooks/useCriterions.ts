import { createCriterion, deleteCriterion, getAllCriterion } from "@api/criterion.api";
import { useMutation, useQuery } from "@tanstack/react-query";
import { ICriterionType } from "types/criterions/criterion.type";

export const useCriterions = () => {
    return useQuery<ICriterionType[]>({
        queryKey: ["criterions"],
        queryFn: () => getAllCriterion(),
    })
}

export const useCreateCriterion = () => {
    return useMutation({
        mutationFn: (name: string) => createCriterion(name),
    })
}

export const useDeleteCriterion = () => {
    return useMutation({
        mutationFn: deleteCriterion,
    })
  }