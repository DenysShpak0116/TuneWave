import { FC, useState } from "react";
import { IconButton, InteractionContainer, InteractionIcon, Logo, LogoContainer } from "./collection-logo.style";
import { AddCriterionModal } from "@modules/AddCriterionToSongModal";
import addCriterionsIcon from "@assets/images/ic_add_criterion.png"
import { ISong } from "types/song/song.type";

interface ICollectionLogoProps {
    logo: string | undefined;
    collectionId: string;
    collectionSongs: ISong[]
}

export const CollectionLogo: FC<ICollectionLogoProps> = ({ logo, collectionSongs, collectionId }) => {
    const [isAddCriterionModalOpen, setIsAddCriterionModalOpen] = useState<boolean>(false);
    return (
        <LogoContainer>
            <Logo src={logo} />
            {collectionSongs && (
                <InteractionContainer>
                    <IconButton onClick={() => setIsAddCriterionModalOpen(true)}>
                        <InteractionIcon src={addCriterionsIcon} />
                    </IconButton>
                </InteractionContainer>
            )}

            <AddCriterionModal
                active={isAddCriterionModalOpen}
                setActive={setIsAddCriterionModalOpen}
                collectionSongs={collectionSongs}
                collectionId={collectionId} />
        </LogoContainer>

    );
};
