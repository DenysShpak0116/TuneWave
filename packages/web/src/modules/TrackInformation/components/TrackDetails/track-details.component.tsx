import { FC } from "react";
import { IAuthor } from "types/song/author.type";
import { ISongTags } from "types/song/songTags.type";
import {
    TrackDetailsContainer,
    TrackInfoText,
    TrackInfoTitle,
    TrackInfoBlock,
    UserInfoText
} from "./track-details.style";
import { ROUTES } from "pages/router/consts/routes.const";
import { parseDate } from "helpers/date-parse";

interface ITrackDetails {
    genre?: string;
    tags?: ISongTags[];
    duration: string;
    date: string;
    artist?: IAuthor[];
    username: string;
    userId: string;
    type: "collection" | "track";
}

export const TrackDetails: FC<ITrackDetails> = ({
    type,
    genre,
    tags,
    duration,
    date,
    artist,
    username,
    userId
}) => {
    const renderBlock = (title: string, value: string | undefined) => (
        <TrackInfoBlock>
            <TrackInfoTitle>{title}</TrackInfoTitle>
            <TrackInfoText>{value || "-"}</TrackInfoText>
        </TrackInfoBlock>
    );

    return (
        <TrackDetailsContainer>
            {type === "track" && (
                <>
                    <TrackInfoBlock>
                        <TrackInfoTitle>Користувач</TrackInfoTitle>
                        <UserInfoText to={ROUTES.USER_PROFILE.replace(':id', userId)}>
                            {username}
                        </UserInfoText>
                    </TrackInfoBlock>
                    {renderBlock("Жанри:", genre)}
                    {renderBlock("Теги:", tags?.map(tag => tag.name).join(", "))}
                    {renderBlock("Тривалість:", duration)}
                    {renderBlock("Дата завантаження:", parseDate(date))}
                    {renderBlock("Головний виконавець:", artist?.map(a => a.name).join(", "))}
                </>
            )}
            {type === "collection" && (
                <>
                    {renderBlock("Тривалість:", duration)}
                    {renderBlock("Дата завантаження:", parseDate(date))}
                    <TrackInfoBlock>
                        <TrackInfoTitle>Користувач</TrackInfoTitle>
                        <UserInfoText to={ROUTES.USER_PROFILE.replace(':id', userId)}>
                            {username}
                        </UserInfoText>
                    </TrackInfoBlock>
                </>
            )}
        </TrackDetailsContainer>
    );
};
