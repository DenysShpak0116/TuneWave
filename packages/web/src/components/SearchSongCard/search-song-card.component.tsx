import { FC } from "react";
import { ISong } from "types/song/song.type";
import { Card, Cover, Icon, Info, Stats, Subtitle, Title } from "./search-song-card.style";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import earIcon from "@assets/images/ic_ear.png"
import likeIcon from "@assets/images/ic_like.png"
import dislikeIcon from "@assets/images/ic_dislike.png"


interface SongCardProps {
    song: ISong;
}

export const SearchSongCard: FC<SongCardProps> = ({ song }) => {
    const navigate = useNavigate();
    return (
        <Card onClick={() => navigate(ROUTES.TRACK_PAGE.replace(":id", song.id))}>
            <Cover src={song.coverUrl} alt={song.title} />
            <Info>
                <Title>{song.title}</Title>
                <Subtitle>
                    {song.user.username} • {song.duration}
                </Subtitle>
                <Stats>
                    <span>{song.genre}</span>
                    <span><Icon src={earIcon} /> {song.listenings}</span>
                    <span><Icon src={likeIcon} /> {song.likes}</span>
                    <span><Icon src={dislikeIcon} /> {song.dislikes}</span>
                </Stats>
            </Info>
        </Card>
    );
};