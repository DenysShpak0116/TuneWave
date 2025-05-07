import { FC, useEffect, useState } from "react";
import { TrackPagePlayerContainer, TrackTitle } from "./track-page-player.style";
import { ISong } from "types/song/song.type";
import { TrackData } from "@components/SongCards/song-cards.component";
import { usePlayerStore } from "@modules/Player/store/player.store";

declare global {
    interface Window {
        Playerjs: any;
    }
}

interface ITrackPagePlayerProps {
    song: ISong;
}

export const TrackPagePlayer: FC<ITrackPagePlayerProps> = ({ song }) => {
    const setTrack = usePlayerStore((state) => state.setTrack);


    useEffect(() => {
        window.PlayerjsEvents = function (event, id, info) {
            if (event === "play") {
                alert("Play triggered");
            }
            if (event === "userplay") {
                console.log("User initiated play", id, info);
            }
            if (event === "time") {
                console.log("Time:", info);
            }
        };

        if (window.Playerjs && song.songUrl) {
            new window.Playerjs({
                id: "trackPagePlayer",
                file: song.songUrl,
                design: 2,
                autoplay: 0,
            });
        }
    }, [song.songUrl]);


    return (
        <TrackPagePlayerContainer>
            <TrackTitle>{song.authors.map(s => s.name).join(',')} - {song.title}</TrackTitle>
            <div id="trackPagePlayer" />
        </TrackPagePlayerContainer>
    );
};
