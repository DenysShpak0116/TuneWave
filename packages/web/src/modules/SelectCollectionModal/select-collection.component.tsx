import { FC, useEffect, useState } from "react";
import { AddCollectionIcon, ModalBody, ModalContent, ModalHeader, ModalHeaderText, Overlay } from "./select-collection.style";
import { CollectionSelectItem, IGetCollectionResponse } from "@components/CollectionSelectItem/collection-select-item.component";
import { useUserCollections } from "./hooks/useUserCollections";
import { Button } from "@ui/Btn/btn.component";
import { useAddSongToCollection } from "./hooks/useAddToCollection";
import addNewCollectionIcon from "@assets/images/ic_add_collection.png"
import { CollectionModal } from "@modules/CollectionModal";

interface ISelectCollection {
    active: boolean;
    setActive: (value: boolean) => void;
    userId: string;
    trackId: string;
}

export const SelectCollectionModal: FC<ISelectCollection> = ({ active, setActive, userId, trackId }) => {
    const [selectedCollection, setSelectedCollection] = useState<IGetCollectionResponse | null>(null)
    const [createCollectionModalActive, setCreateCollectionModalActive] = useState<boolean>(false)
    const { data: collections = [], isLoading, isError, refetch } = useUserCollections(userId);
    const { mutate: addSong } = useAddSongToCollection();

    const handleSubmit = () => {
        if (!selectedCollection) return
        addSong({ trackId: trackId, collectionId: selectedCollection?.id })
    }

    useEffect(() => {
        if (!createCollectionModalActive) {
            refetch();
        }
    }, [createCollectionModalActive]);

    const handleAddNewCollectionClick = () => {
        setCreateCollectionModalActive(true);
    }

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} onClick={(e) => e.stopPropagation()}>
                <ModalHeader>
                    <ModalHeaderText>Оберіть колекцію</ModalHeaderText>
                    <AddCollectionIcon
                        src={addNewCollectionIcon}
                        onClick={handleAddNewCollectionClick} />
                </ModalHeader>
                <ModalBody>
                    {isLoading && <p>Завантаження...</p>}
                    {!isLoading && !isError && collections.map((collection: IGetCollectionResponse) => (
                        <CollectionSelectItem
                            key={collection.id}
                            collection={collection}
                            selected={false}
                            onSelect={() => setSelectedCollection(collection)}
                        />
                    ))}
                </ModalBody>
                <Button text="Додати до колекції" style={{ marginTop: "20px" }} onClick={handleSubmit} />
            </ModalContent>

            <CollectionModal
                active={createCollectionModalActive}
                setActive={setCreateCollectionModalActive}
                userId={userId}
            />
        </Overlay>
    );
};