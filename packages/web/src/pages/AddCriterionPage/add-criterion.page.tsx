import { MainLayout } from "@ui/layout/main-layout";
import { FC, useState } from "react";
import { useCreateCriterion, useCriterions, useDeleteCriterion, useUpdateCriterion } from "./hooks/useCriterions";
import { Button, Container, DeleteButton, Input, List, ListItem, Title, EditButton, SaveButton } from "./add-criterion.style";
import { ICriterionType } from "types/criterions/criterion.type";

export const AddCriterionPage: FC = () => {
    const [name, setName] = useState("");
    const [editValues, setEditValues] = useState<Record<string, string>>({});
    const [editingId, setEditingId] = useState<string | null>(null);

    const { data: criterions, isLoading, refetch } = useCriterions();
    const { mutate: create, isPending: isCreating } = useCreateCriterion();
    const { mutate: remove } = useDeleteCriterion();
    const { mutate: update, isPending: isUpdating } = useUpdateCriterion();

    const handleAddCriterion = () => {
        if (!name.trim()) return;
        create(name, {
            onSuccess: () => {
                setName("");
                refetch();
            },
        });
    };

    const handleDelete = (id: string) => {
        remove(id, {
            onSuccess: () => refetch(),
        });
    };

    const handleStartEdit = (id: string, currentName: string) => {
        setEditingId(id);
        setEditValues({ ...editValues, [id]: currentName });
    };

    const handleSaveEdit = (id: string) => {
        const newName = editValues[id]?.trim();
        if (!newName) return;

        update({ criterionId: id, name: newName }, {
            onSuccess: () => {
                setEditingId(null);
                refetch();
            },
        });
    };

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
                                {editingId === criterion.id ? (
                                    <>
                                        <Input
                                            type="text"
                                            value={editValues[criterion.id] || ""}
                                            onChange={(e) =>
                                                setEditValues({ ...editValues, [criterion.id]: e.target.value })
                                            }
                                        />
                                        <SaveButton onClick={() => handleSaveEdit(criterion.id)} disabled={isUpdating}>
                                            Зберегти
                                        </SaveButton>
                                    </>
                                ) : (
                                    <>
                                        {criterion.name}
                                        <EditButton onClick={() => handleStartEdit(criterion.id, criterion.name)}>
                                            Редагувати
                                        </EditButton>
                                        <DeleteButton onClick={() => handleDelete(criterion.id)}>Видалити</DeleteButton>
                                    </>
                                )}
                            </ListItem>
                        ))}
                    </List>
                )}
            </Container>
        </MainLayout>
    );
};