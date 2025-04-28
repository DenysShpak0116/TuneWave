import { FC, useState } from "react";
import { ISong } from "types/song/song.type";
import { TrackInformationContainer } from "./track-information.style";
import { TrackLogo } from "@modules/TrackInformation/components/TrackLogo/track-logo.component";
import { TrackDetails } from "./components/TrackDetails/track-details.component";
import { parseDate } from "../../helpers/date-parse";
import { CommentSection } from "@components/CommentSection/comment-section.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { IComment } from "types/comments/comment.type";

interface ITrackInformationProps {
    song: ISong;
}

export const TrackInformation: FC<ITrackInformationProps> = ({ song }) => {
    const user = useAuthStore(state => state.user);
    const [comments, setComments] = useState<IComment[]>(song.comments);

    return (
        <TrackInformationContainer>
            <TrackLogo logo={song?.coverUrl} />
            <TrackDetails
                userId={song.user.id}
                username={song.user.username}
                genre={song?.genre}
                tags={song?.songTags}
                duration={song?.duration}
                date={parseDate(song?.createdAt)}
                artist={song?.authors}
            />
            <CommentSection
                userId={user?.id}
                songId={song.id}
                userAvatar={user?.profilePictureUrl}
                comments={comments}
                onNewComment={(comment) => setComments(prev => [...prev, comment])}
                onDeleteComment={(id) => setComments(prev => prev.filter(comment => comment.id !== id))}
            />
        </TrackInformationContainer>
    );
};
