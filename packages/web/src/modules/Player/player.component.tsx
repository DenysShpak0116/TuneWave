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
        Playerjs: any;
    }
}

export const Player: FC = () => {
    const userId = useAuthStore(state => state.user?.id);
    const navigate = useNavigate();

    const {
        trackId,
        trackUrl,
        trackLogo,
        trackName,
        trackArtist,
        shouldAutoPlay,
        setTrack,
        setShouldAutoPlay,
        setIsPlaying,
    } = usePlayerStore();

    const [isModalOpen, setIsModalOpen] = useState(false);

    useEffect(() => {
        if (!trackUrl) {
            const savedTrack = localStorage.getItem("current-track");
            if (savedTrack) {
                const track = JSON.parse(savedTrack);
                setTrack(track);
            }
            return;
        }

        const player = new window.Playerjs({
            id: "player",
            file: trackUrl,
            autoplay: shouldAutoPlay ? 1 : 0,
        });

        if (shouldAutoPlay) {
            setShouldAutoPlay(false);
        }

        window.PlayerjsEvents = function (event: string, id: string, info: any) {
            if (id !== "player") return;

            if (event === "time") {
                localStorage.setItem("player-timeline", JSON.stringify({ time: info }));
            }

            if (event === "play" || event === "userplay") {
                localStorage.setItem("player-sync", JSON.stringify({ type: "pauseOthers", excludeId: "player" }));
                setIsPlaying(true);
            }

            if (event === "pause") {
                setIsPlaying(false);
            }
        };

        localStorage.setItem("current-track", JSON.stringify({
            trackId,
            trackUrl,
            trackLogo,
            trackName,
            trackArtist
        }));

        const syncHandler = (e: StorageEvent) => {
            if (e.key === "player-sync" && e.newValue) {
                const { type, value, excludeId } = JSON.parse(e.newValue);
                if (type === "play") player.api("play");
                if (type === "pause") player.api("pause");
                if (type === "seek") player.api("seek", value);
                if (type === "volume") player.api("volume", value);
                if (type === "pauseOthers" && excludeId !== "player") {
                    player.api("pause");
                }
            }
        };
        window.addEventListener("storage", syncHandler);

        return () => {
            window.removeEventListener("storage", syncHandler);
        };
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

            <SelectCollectionModal
                trackId={trackId}
                userId={userId!}
                active={isModalOpen}
                setActive={setIsModalOpen}
            />
        </>
    );
};
