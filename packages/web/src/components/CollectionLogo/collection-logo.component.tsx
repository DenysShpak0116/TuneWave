import { FC, useState } from "react";
import { IconButton, InteractionContainer, InteractionIcon, Logo, LogoContainer } from "./collection-logo.style";
import { AddCriterionModal } from "@modules/AddCriterionToSongModal";
import addCriterionsIcon from "@assets/images/ic_add_criterion.png"
import { ISong } from "types/song/song.type";
import stepsIcon from "@assets/images/ic_steps.png"
import { RateSongsModal } from "@modules/RateSongsModal/rate-songs-modal.component";
import crossIcon from "@assets/images/ic_cross.png"
import { ConfirmDeleteModal } from "@components/ConfirmDeleteModal/confirmDelete.modal";
import { useNavigate } from "react-router-dom";
import { useDeleteCollection } from "./hooks/useDeleteCollection";
import { ROUTES } from "pages/router/consts/routes.const";
import plusIcon from "@assets/images/ic_plus.png"
import { useAddCollectionToUser } from "./hooks/useAddCollectionToUser";


interface ICollectionLogoProps {
    logo: string | undefined;
    collectionId: string;
    collectionSongs: ISong[];
    hasAllVectors: boolean;
    isMainUserCollection: boolean
}

export const CollectionLogo: FC<ICollectionLogoProps> = ({ logo, collectionSongs, collectionId, hasAllVectors, isMainUserCollection }) => {
    const navigate = useNavigate()
    const [isAddCriterionModalOpen, setIsAddCriterionModalOpen] = useState<boolean>(false);
    const [isRateModalOpen, setIsRateModalOpen] = useState<boolean>(false)
    const [isDeleteConfirmationModalOpen, setIsDeleteConfirmationModalOpen] = useState<boolean>(false);
    const { mutate: deleteCollectionMutate } = useDeleteCollection();
    const { mutate: addCollectionMutate } = useAddCollectionToUser()

    const handleDelete = () => {
        deleteCollectionMutate(collectionId)
        navigate(ROUTES.HOME)
    }

    const handleAdd = () => {
        addCollectionMutate(collectionId)
    }

    return (
        <LogoContainer>
            <Logo src={logo} />
            {collectionSongs && (
                <InteractionContainer>
                    <IconButton onClick={() => setIsAddCriterionModalOpen(true)}>
                        <InteractionIcon src={addCriterionsIcon} />
                    </IconButton>
                    {hasAllVectors && (
                        <IconButton onClick={() => setIsRateModalOpen(true)}>
                            <InteractionIcon src={stepsIcon} />
                        </IconButton>
                    )}
                    {isMainUserCollection && (
                        <IconButton onClick={() => setIsDeleteConfirmationModalOpen(true)}>
                            <InteractionIcon src={crossIcon} />
                        </IconButton>
                    )}

                    {!isMainUserCollection && (
                        <IconButton onClick={handleAdd}>
                            <InteractionIcon src={plusIcon} />
                        </IconButton>
                    )}

                </InteractionContainer>
            )}

            <RateSongsModal
                active={isRateModalOpen}
                setActive={setIsRateModalOpen}
                collectionSongs={collectionSongs}
                collectionId={collectionId}
            />

            <AddCriterionModal
                active={isAddCriterionModalOpen}
                setActive={setIsAddCriterionModalOpen}
                collectionSongs={collectionSongs}
                collectionId={collectionId}
            />

            <ConfirmDeleteModal
                text="Ви впевнені, що хочете видалити колекцію?"
                active={isDeleteConfirmationModalOpen}
                setActive={setIsDeleteConfirmationModalOpen}
                onDelete={handleDelete}
            />
        </LogoContainer>

    );
};
