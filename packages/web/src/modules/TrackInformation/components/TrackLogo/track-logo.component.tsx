import { FC, useState } from "react";
import { InteractionContainer, Logo, LogoContainer, IconButton, InteractionIcon } from "./track-logo.style";
import likeIcon from '@assets/images/ic_like.png'
import likeFilledIcon from '@assets/images/ic_like_filled.png'
import dislikeIcon from '@assets/images/ic_dislike.png'
import dislikeFilledIcon from '@assets/images/ic_dislike_filled.png'

type ReactionType = "like" | "dislike" | "none";

interface ITrackLogo {
    logo: string | undefined;
    reactFn: (type: "like" | "dislike") => void
    type: { type: ReactionType }
}

export const TrackLogo: FC<ITrackLogo> = ({ logo, reactFn, type: { type } }) => {
    const [reaction, setReaction] = useState<"like" | "dislike" | null>(type === 'none' ? null : type);



    const handleReaction = (type: "like" | "dislike") => {
        const isSameReaction = reaction === type;
        const newReaction = isSameReaction ? null : type;

        setReaction(newReaction);
        reactFn(type);
    };

    return (
        <LogoContainer>
            <Logo src={logo} />
            <InteractionContainer>
                <IconButton onClick={() => handleReaction("like")}>
                    <InteractionIcon src={reaction === "like" ? likeFilledIcon : likeIcon} />
                </IconButton>
                <IconButton onClick={() => handleReaction("dislike")}>
                    <InteractionIcon src={reaction === "dislike" ? dislikeFilledIcon : dislikeIcon} />
                </IconButton>
            </InteractionContainer>
        </LogoContainer>
    );
};
