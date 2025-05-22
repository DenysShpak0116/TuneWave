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
import { Button } from "@ui/Btn/btn.component";
import { useNavigate } from "react-router-dom";

interface ITrackDetails {
    trackId?: string;
    collectionId?: string;
    genre?: string;
    tags?: ISongTags[];
    duration: string;
    date: string;
    artist?: IAuthor[];
    username: string;
    userId: string;
    type: "collection" | "track";
    isMainUser: boolean,
    collectionName?: string,
    collectionDescription?: string;
}

export const TrackDetails: FC<ITrackDetails> = ({
    trackId,
    collectionId,
    isMainUser,
    type,
    genre,
    tags,
    duration,
    date,
    artist,
    username,
    userId,
    collectionName,
    collectionDescription,
}) => {
    const navigate = useNavigate();

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
                    {isMainUser && (
                        <TrackInfoBlock>
                            <Button
                                onClick={() => navigate(ROUTES.UPDATE_TRACK_PAGE.replace(":id", trackId!))}
                                text="Редагувати"
                                style={{ padding: "5px 10px", fontSize: "14px" }} />
                        </TrackInfoBlock>
                    )}
                    <TrackInfoBlock>
                        <TrackInfoTitle>Користувач</TrackInfoTitle>
                        <UserInfoText to={ROUTES.USER_PROFILE.replace(':id', userId)}>
                            {username}
                        </UserInfoText>
                    </TrackInfoBlock>

                    {renderBlock("Жанри:", genre)}
                    {renderBlock("Теги:", tags?.map(tag => tag.name).join(", "))}
                    {renderBlock("Тривалість:", duration)}
                    {renderBlock("Дата завантаження:", date)}
                    {renderBlock("Головний виконавець:", artist?.map(a => a.name).join(", "))}
                </>
            )}
            {type === "collection" && (
                <>
                    <TrackInfoBlock>
                        <Button
                            onClick={() => navigate(ROUTES.SONGS_CRITERIONS_PAGE.replace(":id", collectionId!))}
                            text="Таблиця крітерієв"
                            style={{ padding: "5px 10px", fontSize: "14px" }} />
                    </TrackInfoBlock>
                    <TrackInfoBlock>
                        <Button
                            onClick={() => navigate(ROUTES.COLLECTIVE_DECISION_PAGE.replace(":id", collectionId!))}
                            text="Результати голусування"
                            style={{ padding: "5px 10px", fontSize: "14px" }} />
                    </TrackInfoBlock>
                    <TrackInfoBlock>
                        <Button
                            onClick={() => navigate(ROUTES.UPDATE_COLLECTION.replace(":id", collectionId!))}
                            text="Редагувати"
                            style={{ padding: "5px 10px", fontSize: "14px" }} />
                    </TrackInfoBlock>
                    {renderBlock("Назва колекції:", collectionName)}
                    {renderBlock("Опис колекції:", collectionDescription)}
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
