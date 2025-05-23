import { FC, useState } from "react";
import { IconButton, InteractionContainer, InteractionIcon, Logo, LogoContainer } from "./collection-logo.style";
import { AddCriterionModal } from "@modules/AddCriterionToSongModal";
import addCriterionsIcon from "@assets/images/ic_add_criterion.png"
import { ISong } from "types/song/song.type";
import stepsIcon from "@assets/images/ic_steps.png"
import { RateSongsModal } from "@modules/RateSongsModal/rate-songs-modal.component";

interface ICollectionLogoProps {
    logo: string | undefined;
    collectionId: string;
    collectionSongs: ISong[];
    hasAllVectors: boolean;
}

export const CollectionLogo: FC<ICollectionLogoProps> = ({ logo, collectionSongs, collectionId, hasAllVectors }) => {
    const [isAddCriterionModalOpen, setIsAddCriterionModalOpen] = useState<boolean>(false);
    const [isRateModalOpen, setIsRateModalOpen] = useState<boolean>(false)
    
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
        </LogoContainer>

    );
};
