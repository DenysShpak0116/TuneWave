import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";
import { useCollectiveDecision } from "./hooks/useCollectiveDecision";
import { Loader } from "@ui/Loader/loader.component";
import {
    PageWrapper,
    SectionTitle,
    SongList,
    SongItem,
    ProfileTable,
    TableRow,
    TableCell,
    UserName,
    SongListCell,
} from "./collective-decision.style";

export const CollectiveDecisionPage: FC = () => {
    const { id } = useParams();
    const { data, isLoading } = useCollectiveDecision(id!);

    if (isLoading || !data) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    const { collectiveRank, profileTable } = data;

    return (
        <MainLayout>
            <PageWrapper>
                <SectionTitle>Топ пісень:</SectionTitle>
                <SongList>
                    {Object.entries(collectiveRank).map(([rank, { songName }]) => (
                        <SongItem key={rank}>
                            {rank}. {songName}
                        </SongItem>
                    ))}
                </SongList>

                <SectionTitle>Профілі користувачів:</SectionTitle>
                <ProfileTable>
                    <thead>
                        <TableRow>
                            <TableCell>Ранг</TableCell>
                            {Object.values(profileTable)[0] &&
                                Object.keys(Object.values(profileTable)[0]).map((user) => (
                                    <TableCell key={user}>
                                        <UserName>{user.split(":")[1]}</UserName>
                                    </TableCell>
                                ))}
                        </TableRow>
                    </thead>
                    <tbody>
                        {Object.entries(profileTable).map(([rank, userMap]) => (
                            <TableRow key={rank}>
                                <TableCell>{rank}</TableCell>
                                {Object.values(userMap).map((songs, idx) => (
                                    <SongListCell key={idx}>
                                        {songs.join(", ")}
                                    </SongListCell>
                                ))}
                            </TableRow>
                        ))}
                    </tbody>
                </ProfileTable>
            </PageWrapper>
        </MainLayout>
    );
};