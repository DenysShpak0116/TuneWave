import { FC, useMemo } from "react";
import { useGetTopTracks, useGetTracksWithPopular } from "./hooks/useGetTopTracks";
import * as S from "./top-listens.style";
import { Loader } from "@ui/Loader/loader.component";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { Bar, BarChart, ResponsiveContainer, Tooltip, XAxis } from "recharts";
import { ITopTrackStatistic } from "types/song/song-statistic";

export const TopListens: FC = () => {
    const { data: top5Data, isLoading: isTop5Loading } = useGetTopTracks();
    const { data: popularData, isLoading: isPopularLoading } = useGetTracksWithPopular();
    const navigate = useNavigate();
    const isLoading = isTop5Loading || isPopularLoading;

    const topTracks = useMemo(() => top5Data?.slice(0, 5) || [], [top5Data]);

    const chartData = useMemo(
        () =>
            topTracks.map((item, index) => ({
                name: `#${index + 1}`,
                listens: item.count,
            })),
        [topTracks]
    );

    const mostPopular = useMemo(() => {
        if (!popularData || popularData.length === 0) return null;
        return popularData.reduce(
            (prev: ITopTrackStatistic, curr: ITopTrackStatistic) => (curr.song.listenings > prev.song.listenings ? curr : prev),
            popularData[0]
        );
    }, [popularData]);

    if (isLoading) {
        return (
            <S.LoaderWrapper>
                <Loader />
            </S.LoaderWrapper>
        );
    }

    if (topTracks.length === 0) {
        return <S.Empty>Немає даних для відображення</S.Empty>;
    }

    return (
        <S.Container>
            <S.Title>Топ 5 треків</S.Title>

            <S.Content>
                <S.LeftColumn>
                    <S.TrackList>
                        {topTracks.map((item, index) => (
                            <S.TrackCard
                                key={item.song.id}
                                rank={index + 1}
                                onClick={() => navigate(ROUTES.TRACK_PAGE.replace(":id", item.song.id))}
                            >
                                <S.RankBadge rank={index + 1}>#{index + 1}</S.RankBadge>
                                <S.CoverWrapper>
                                    <S.Cover src={item.song.coverUrl} alt={item.song.title} />
                                    <S.Overlay />
                                </S.CoverWrapper>
                                <S.Info>
                                    <S.TrackTitle>{item.song.title}</S.TrackTitle>
                                    <S.Genre>{item.song.genre}</S.Genre>
                                    <S.Duration>{item.song.duration}</S.Duration>
                                </S.Info>
                                <S.Count>{item.count} прослуховувань</S.Count>
                            </S.TrackCard>
                        ))}
                    </S.TrackList>

                    <S.ChartCard>
                        <h3>Порівняння кількості прослуховувань</h3>
                        <S.ChartWrapper>
                            <ResponsiveContainer width="100%" height={250}>
                                <BarChart data={chartData}>
                                    <XAxis dataKey="name" stroke="#ccc" />
                                    <Tooltip />
                                    <Bar dataKey="listens" radius={[6, 6, 0, 0]} fill="#4A90E2" />
                                </BarChart>
                            </ResponsiveContainer>
                        </S.ChartWrapper>
                    </S.ChartCard>
                </S.LeftColumn>

                {popularData && popularData.length > 0 && (
                    <S.PopularContainer>
                        <S.Title>Часто слухані треки з найпопулярним</S.Title>
                        <S.PopularList>
                            {popularData.map((item: ITopTrackStatistic) => {
                                const isTop = item.song.id === mostPopular.song.id;
                                return (
                                    <S.PopularCard
                                        key={item.song.id}
                                        isTop={isTop}
                                        onClick={() => navigate(ROUTES.TRACK_PAGE.replace(":id", item.song.id))}
                                    >
                                        <S.CoverWrapper>
                                            <S.Cover src={item.song.coverUrl} alt={item.song.title} />
                                        </S.CoverWrapper>
                                        <S.Info>
                                            <S.TrackTitle>{item.song.title}</S.TrackTitle>
                                            <S.Genre>{item.song.genre}</S.Genre>
                                            <S.Duration>{item.song.duration}</S.Duration>
                                        </S.Info>
                                        <S.Count>{item.count} разів</S.Count>
                                    </S.PopularCard>
                                );
                            })}
                        </S.PopularList>
                    </S.PopularContainer>
                )}
            </S.Content>
        </S.Container>
    );
};
