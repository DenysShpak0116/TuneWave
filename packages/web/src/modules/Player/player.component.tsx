import { FC, useEffect } from "react";
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

declare global {
    interface Window {
        Playerjs: any;
    }
}

export const Player: FC = () => {
    const navigate = useNavigate();
    const { trackId, trackUrl, trackLogo, trackName, trackArtist } = usePlayerStore();
    console.log(trackUrl);


    useEffect(() => {
        if (trackUrl && window.Playerjs) {
            new window.Playerjs({
                id: "player",
                file: trackUrl,
                autoplay: 1,

            });
        }
    }, [trackUrl]);

    if (!trackUrl) return null;

    return (
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
            <AddIcon src={addToCollection} />
        </PlayerContainer>
    );
};