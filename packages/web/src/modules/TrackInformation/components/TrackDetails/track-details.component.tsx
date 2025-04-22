import { FC } from "react";
import { IAuthor } from "types/song/author.type";
import { ISongTags } from "types/song/songTags.type";
import {
    TrackDetailsContainer,
    TrackInfoText,
    TrackInfoTitle,
    TrackInfoBlock
} from "./track-details.style";

interface ITrackDetails {
    genre: string;
    tags: ISongTags[];
    duration: string;
    date: string;
    artist: IAuthor[];
}

export const TrackDetails: FC<ITrackDetails> = ({ genre, tags, duration, date, artist }) => {
    return (
        <TrackDetailsContainer>
            <TrackInfoBlock>
                <TrackInfoTitle>Жанри:</TrackInfoTitle>
                <TrackInfoText>{genre}</TrackInfoText>
            </TrackInfoBlock>
            <TrackInfoBlock>
                <TrackInfoTitle>Теги:</TrackInfoTitle>
                <TrackInfoText>{tags.map(tag => tag.name).join(', ')}</TrackInfoText>
            </TrackInfoBlock>
            <TrackInfoBlock>
                <TrackInfoTitle>Тривалість:</TrackInfoTitle>
                <TrackInfoText>{duration}</TrackInfoText>
            </TrackInfoBlock>
            <TrackInfoBlock>
                <TrackInfoTitle>Дата завантаження:</TrackInfoTitle>
                <TrackInfoText>{date}</TrackInfoText>
            </TrackInfoBlock>
            <TrackInfoBlock>
                <TrackInfoTitle>Головний виконавець:</TrackInfoTitle>
                <TrackInfoText>{artist.map(a => a.name).join(', ')}</TrackInfoText>
            </TrackInfoBlock>
        </TrackDetailsContainer>
    );
};
