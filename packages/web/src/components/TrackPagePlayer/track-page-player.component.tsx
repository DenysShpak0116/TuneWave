import { FC, useEffect } from "react";
import { ListeningContainer, TrackPagePlayerContainer, TrackTitle } from "./track-page-player.style";
import { ISong } from "types/song/song.type";
import { usePlayerStore } from "@modules/Player/store/player.store";
import musicNote from "@assets/images/ic_musical_note.png"

interface ITrackPagePlayerProps {
    song: ISong;
    songs?: ISong[];
}

export const TrackPagePlayer: FC<ITrackPagePlayerProps> = ({ song, songs }) => {
    const setTrack = usePlayerStore((state) => state.setTrack);

    useEffect(() => {
        if (!window.Playerjs || !song.songUrl) return;

        const player = new window.Playerjs({
            id: "trackPagePlayer",
            file: song.songUrl,
            design: 2,
            autoplay: 0,
        });


        window.PlayerjsEvents = function (event: string, id: string, info: any) {
            if (event === "play" || event === "userplay") {
                setTrack({
                    trackId: song.id,
                    trackUrl: song.songUrl,
                    trackName: song.title,
                    trackLogo: song.coverUrl,
                    trackArtist: song.user.username
                });

                localStorage.setItem("player-sync", JSON.stringify({
                    type: "pauseOthers",
                    excludeId: "trackPagePlayer"
                }));
            }
        };

        const syncHandler = (e: StorageEvent) => {
            if (e.key === "player-sync" && e.newValue) {
                const { type, value, excludeId } = JSON.parse(e.newValue);
                if (excludeId === "trackPagePlayer") return;

                if (type === "pauseOthers") player.api("pause");
                if (type === "seek") player.api("seek", value);
                if (type === "play") player.api("play");
                if (type === "pause") player.api("pause");
                if (type === "volume") player.api("volume", value);
            }
        };

        const timeSyncHandler = (e: StorageEvent) => {
            if (e.key === "player-timeline" && e.newValue) {
                const { time } = JSON.parse(e.newValue);
                player.api("seek", time);
            }
        };
        window.addEventListener("storage", timeSyncHandler);

        return () => {
            window.removeEventListener("storage", syncHandler);
            window.removeEventListener("storage", timeSyncHandler);
        };
    }, [song.songUrl]);



    return (
        <TrackPagePlayerContainer>
            <ListeningContainer>
                <img src={musicNote} alt="note" />
                <p>{song.listenings}</p>
            </ListeningContainer>
            <TrackTitle>
                {song.authors?.map((s) => s.name).join(", ") ?? "Невідомий автор"} - {song.title}
            </TrackTitle>
            <div id="trackPagePlayer" />
        </TrackPagePlayerContainer>
    );
};
