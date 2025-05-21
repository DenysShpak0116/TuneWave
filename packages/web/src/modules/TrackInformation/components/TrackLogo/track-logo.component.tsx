import { FC, useState } from "react";
import { InteractionContainer, Logo, LogoContainer, IconButton, InteractionIcon } from "./track-logo.style";
import likeIcon from '@assets/images/ic_like.png';
import likeFilledIcon from '@assets/images/ic_like_filled.png';
import dislikeIcon from '@assets/images/ic_dislike.png';
import dislikeFilledIcon from '@assets/images/ic_dislike_filled.png';
import 'react-medium-image-zoom/dist/styles.css';
import crossIcon from "@assets/images/ic_cross.png"
import { useDeleteTrack } from "pages/TrackPage/hooks/useGetTrack";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import plusIcon from "@assets/images/ic_plus.png"
import { SelectCollectionModal } from "@modules/SelectCollectionModal";

type ReactionType = "like" | "dislike" | "none";

interface ITrackLogo {
    songId: string;
    logo: string | undefined;
    userId: string | undefined;
    reactFn: (type: "like" | "dislike") => void;
    type: { type: ReactionType };
    isUserMainTrack: boolean
}

export const TrackLogo: FC<ITrackLogo> = ({ logo, reactFn, type: { type }, userId, isUserMainTrack, songId }) => {
    const navigate = useNavigate();
    const [reaction, setReaction] = useState<"like" | "dislike" | null>(type === 'none' ? null : type);
    const { mutate: deleteTrackMutate } = useDeleteTrack();
    const [isAddToCollectionModalOpen, setIsAddToCollectionModalOpen] = useState<boolean>(false);

    const handleReaction = (type: "like" | "dislike") => {
        const isSameReaction = reaction === type;
        const newReaction = isSameReaction ? null : type;

        setReaction(newReaction);
        reactFn(type);
    };

    const handleDelete = () => {
        deleteTrackMutate(songId);
        navigate(ROUTES.HOME)
    };

    return (
        <LogoContainer>
            <Logo src={logo} alt="Track logo" />
            {userId && (
                <InteractionContainer>
                    <IconButton onClick={() => handleReaction("like")}>
                        <InteractionIcon src={reaction === "like" ? likeFilledIcon : likeIcon} />
                    </IconButton>
                    <IconButton onClick={() => handleReaction("dislike")}>
                        <InteractionIcon src={reaction === "dislike" ? dislikeFilledIcon : dislikeIcon} />
                    </IconButton>
                    <IconButton onClick={() => setIsAddToCollectionModalOpen(true)}>
                        <InteractionIcon src={plusIcon} />
                    </IconButton>
                    {isUserMainTrack && (
                        <IconButton onClick={handleDelete}>
                            <InteractionIcon src={crossIcon} />
                        </IconButton>
                    )}
                </InteractionContainer>
            )}

            <SelectCollectionModal
                trackId={songId}
                userId={userId!}
                active={isAddToCollectionModalOpen}
                setActive={setIsAddToCollectionModalOpen}
            />
        </LogoContainer>
    );
};
