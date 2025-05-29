import { FC, useEffect, useRef, useState } from "react";
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
import { useAddListening } from "./hooks/useAddListening";

declare global {
    interface Window {
        Playerjs: any;
        PlayerjsEvents: any;
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
        playlist,
        isPlaying,
    } = usePlayerStore();

    const [isModalOpen, setIsModalOpen] = useState(false);
    const { mutate: addListening } = useAddListening();
    const playerRef = useRef<any>(null);
    const prevTrackRef = useRef<string | null>(null);

    const mappedPlaylist = playlist.length > 0
        ? playlist.map(song => ({ title: song.title, file: song.file }))
        : [];

    useEffect(() => {
        if (!trackUrl) return;

        if (!playerRef.current) {
            playerRef.current = new window.Playerjs({
                id: "player",
                playlist: mappedPlaylist.length > 0 ? mappedPlaylist : [{ title: trackName, file: trackUrl }],
                autoplay: shouldAutoPlay ? 1 : 0,
            });

            window.PlayerjsEvents = (event: string, id: string, info: any) => {
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

                if (event === "ended") {
                    const currentIndex = playlist.findIndex(p => p.file === trackUrl);
                    const nextTrack = playlist[currentIndex + 1];
                    if (nextTrack) {
                        setTrack({
                            trackUrl: nextTrack.file,
                            trackName: nextTrack.title,
                            trackArtist: nextTrack.artist || "",
                            trackId: nextTrack.id,
                            trackLogo: nextTrack.logo || "",
                        });
                    }
                }

                if (event === "playlist") {
                    const index = info;
                    const currentTrack = playlist[index];
                    if (currentTrack) {
                        setTrack({
                            trackUrl: currentTrack.file,
                            trackName: currentTrack.title,
                            trackArtist: currentTrack.artist || "",
                            trackId: currentTrack.id,
                            trackLogo: currentTrack.logo || "",
                        });
                    }
                }
            };

            const syncHandler = (e: StorageEvent) => {
                if (e.key === "player-sync" && e.newValue) {
                    const { type, value, excludeId } = JSON.parse(e.newValue);
                    if (!playerRef.current) return;
                    if (type === "play") playerRef.current.api("play");
                    if (type === "pause") playerRef.current.api("pause");
                    if (type === "seek") playerRef.current.api("seek", value);
                    if (type === "volume") playerRef.current.api("volume", value);
                    if (type === "pauseOthers" && excludeId !== "player") {
                        playerRef.current.api("pause");
                    }
                }
            };

            window.addEventListener("storage", syncHandler);

            return () => {
                window.removeEventListener("storage", syncHandler);
            };
        }
    }, [trackUrl]);

    useEffect(() => {
        if (!playerRef.current || !trackUrl) return;

        const player = playerRef.current;
        const isNewTrack = prevTrackRef.current !== trackUrl;
        const playerTimeLine = isNewTrack ? 0 : player.api("time");

        if (playlist.length > 0) {
            const mapped = playlist.map(song => ({
                title: song.title,
                file: song.file,
            }));
            player.api("playlist", mapped);
        } else {
            player.api("playlist", [{ title: trackName, file: trackUrl }]);
        }

        if (isNewTrack) {
            player.api("file", trackUrl);
            prevTrackRef.current = trackUrl;
        }

        player.api("seek", playerTimeLine);

        if (shouldAutoPlay) {
            player.api("play");
            setShouldAutoPlay(false);
        } else if (!isPlaying) {
            player.api("pause");
        }

        localStorage.setItem("current-track", JSON.stringify({
            trackId,
            trackUrl,
            trackLogo,
            trackName,
            trackArtist
        }));

        if (userId) {
            addListening({ songId: trackId, userId });
        }

        const interval = setInterval(() => {
            const playlistIdRaw = player.api("playlist_id");
            const currentIndex = Number(playlistIdRaw?.split("-")[1]);
            const currentTrack = playlist[currentIndex];

            if (currentTrack && currentTrack.file !== prevTrackRef.current) {
                setTrack({
                    trackUrl: currentTrack.file,
                    trackName: currentTrack.title,
                    trackArtist: currentTrack.artist || "",
                    trackId: currentTrack.id,
                    trackLogo: currentTrack.logo || "",
                });
                prevTrackRef.current = currentTrack.file;
            }
        }, 500);

        return () => clearInterval(interval);
    }, [trackUrl, isPlaying, playlist]);

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
