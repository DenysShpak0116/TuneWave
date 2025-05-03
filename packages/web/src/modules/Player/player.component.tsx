import { FC, useEffect, useState } from "react";
import { PlayerContainer } from "@components/PlayerContainer/player-container.component";
import {
    TrackLogo,
    TrackInfoWrapper,
    TrackName,
    TrackArtist,
    TextBlock,
    AddIcon
} from "./player.style";
import addToCollection from "@assets/images/ic_add_to_collection.png";
import { usePlayerStore } from "./store/player.store";
import { ROUTES } from "pages/router/consts/routes.const";
import { useNavigate } from "react-router-dom";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { SelectCollectionModal } from "@modules/SelectCollectionModal";

declare global {
    interface Window {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        Playerjs: any;
    }
}

export const Player: FC = () => {
    const userId = useAuthStore(state => state.user?.id)
    const navigate = useNavigate();
    const { trackId, trackUrl, trackLogo, trackName, trackArtist } = usePlayerStore();
    const [isModalOpen, setIsModalOpen] = useState(false);

    useEffect(() => {
        if (trackUrl && window.Playerjs) {
            const player = new window.Playerjs({
                id: "player",
                file: trackUrl,
                autoplay: 0,

            });
            player.api("play", trackUrl)
        }
    }, [trackUrl]);


    if (!trackUrl) return null;

    return (
        <>
            <PlayerContainer>
                <TrackInfoWrapper>
                    <TrackLogo
                        onClick={() => {
                            navigate(ROUTES.TRACK_PAGE.replace(":id", trackId));
                        }}
                        src={trackLogo}
                        alt={trackName}
                    />
                    <TextBlock>
                        <TrackName>{trackName}</TrackName>
                        <TrackArtist>{trackArtist}</TrackArtist>
                    </TextBlock>
                </TrackInfoWrapper>

                <div style={{ width: "700px", height: "35px" }}>
                    <div id="player" />
                </div>
                {userId && (
                    <AddIcon src={addToCollection} onClick={() => setIsModalOpen(true)} />
                )}
            </PlayerContainer>
            <SelectCollectionModal trackId={trackId} userId={userId!} active={isModalOpen} setActive={setIsModalOpen} />
        </>
    );
};