import { MainLayout } from "@ui/layout/main-layout";
import { FC, useState } from "react";
import { useCreateCriterion, useCriterions, useDeleteCriterion } from "./hooks/useCriterions";
import { Button, Container, DeleteButton, Input, List, ListItem, Title } from "./add-criterion.style";
import { ICriterionType } from "types/criterions/criterion.type";

export const AddCriterionPage: FC = () => {
    const [name, setName] = useState("")
    const { data: criterions, isLoading, refetch } = useCriterions()
    const { mutate: create, isPending: isCreating } = useCreateCriterion()
    const { mutate: remove } = useDeleteCriterion()

    const handleAddCriterion = () => {
        if (!name.trim()) return
        create(name, {
            onSuccess: () => {
                setName("")
                refetch()
            },
        })
    }

    const handleDelete = (id: string) => {
        remove(id, {
            onSuccess: () => refetch(),
        })
    }

    return (
        <MainLayout>
            <Container>
                <Title>Критерії</Title>

                <Input
                    type="text"
                    placeholder="Назва критерію"
                    value={name}
                    onChange={(e) => setName(e.target.value)}
                />

                <Button onClick={handleAddCriterion} disabled={isCreating || !name.trim()}>
                    Додати критерій
                </Button>

                {isLoading ? (
                    <p>Завантаження...</p>
                ) : (
                    <List>
                        {criterions?.map((criterion: ICriterionType) => (
                            <ListItem key={criterion.id}>
                                {criterion.name}
                                <DeleteButton onClick={() => handleDelete(criterion.id)}>Видалити</DeleteButton>
                            </ListItem>
                        ))}
                    </List>
                )}
            </Container>
        </MainLayout>
    )
}