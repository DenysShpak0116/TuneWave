import { FC } from "react";
import { ISong } from "types/song/song.type";
import { TrackInformationContainer } from "./track-information.style";
import { TrackLogo } from "@modules/TrackInformation/components/TrackLogo/track-logo.component";
import { TrackDetails } from "./components/TrackDetails/track-details.component";
import { parseDate } from "./helpers/date-parse";

interface ITrackInformationProps {
    song: ISong | undefined
}

export const TrackInformation: FC<ITrackInformationProps> = ({ song }) => {

    return (
        <TrackInformationContainer>
            <TrackLogo logo={song?.coverUrl} />
            <TrackDetails
                genre={song?.genre}
                tags={song?.songTags}
                duration={song?.duration}
                date={parseDate(song?.createdAt)}
                artist={song?.authors} />
        </TrackInformationContainer>
    )
}