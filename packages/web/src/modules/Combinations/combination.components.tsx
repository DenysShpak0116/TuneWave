import { FC, useMemo, useState } from "react";
import * as S from "./combination.style";
import { useGetCommonCombinations, useGetRareCombinations } from "./hooks/useGetCombinations";
import { Loader } from "@ui/Loader/loader.component";
import { useGetTrack } from "pages/TrackPage/hooks/useGetTrack";

export const Combination: FC = () => {
    const [activeTab, setActiveTab] = useState<"common" | "rare">("common");

    const common = useGetCommonCombinations();
    const rare = useGetRareCombinations();

    const data = activeTab === "common" ? common.data : rare.data;
    const isLoading = common.isLoading || rare.isLoading;

    const chartData = useMemo(() => {
        if (!data) return [];
        return data.map((item: any) => ({
            id: item.Key,
            listens: item.Value,
        }));
    }, [data]);

    if (isLoading) {
        return (
            <S.LoaderWrapper>
                <Loader />
            </S.LoaderWrapper>
        );
    }

    return (
        <S.Container>
            <S.Title>Комбінації треків</S.Title>

            <S.Tabs>
                <S.TabButton active={activeTab === "common"} onClick={() => setActiveTab("common")}>
                    Часті
                </S.TabButton>
                <S.TabButton active={activeTab === "rare"} onClick={() => setActiveTab("rare")}>
                    Рідкісні
                </S.TabButton>
            </S.Tabs>

            {chartData.length === 0 ? (
                <S.Empty>Даних немає</S.Empty>
            ) : (
                <>
                    {/* <S.ChartCard>
                        <h3>Кількість прослуховувань комбінацій</h3>
                        <ResponsiveContainer width="100%" height={300}>
                            <BarChart data={chartData}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="id" hide />
                                <YAxis />
                                <Tooltip />
                                <Bar dataKey="listens" fill="#4A90E2" radius={[6, 6, 0, 0]} />
                            </BarChart>
                        </ResponsiveContainer>
                    </S.ChartCard> */}

                    <S.CombinationList>
                        {chartData.slice(0, 6).map(({ id, listens }) => (
                            <CombinationCard key={id} id={id} listens={listens} />
                        ))}
                    </S.CombinationList>
                </>
            )}
        </S.Container>
    );
};

const CombinationCard: FC<{ id: string; listens: number }> = ({ id, listens }) => {
    const trackIds = id.split(",").filter(Boolean);
    const trackCountMap = trackIds.reduce<Record<string, number>>((acc, trackId) => {
        acc[trackId] = (acc[trackId] || 0) + 1;
        return acc;
    }, {});
    const uniqueTracks = Object.keys(trackCountMap);

    return (
        <S.ComboCard>
            <S.CoversRow>
                {uniqueTracks.slice(0, 4).map((trackId) => (
                    <TrackPreview key={trackId} id={trackId} count={trackCountMap[trackId]} />
                ))}
            </S.CoversRow>
            <S.ComboFooter>
                <S.ComboText>{uniqueTracks.length} треків</S.ComboText>
                <S.ComboCount>Використано {listens} разів</S.ComboCount>
            </S.ComboFooter>
        </S.ComboCard>
    );
};

const TrackPreview: FC<{ id: string; count: number }> = ({ id, count }) => {
    const { data: track } = useGetTrack(id);
    if (!track) return <S.TrackCoverSkeleton />;

    return (
        <S.TrackPreviewWrapper>
            <S.TrackCover src={track.coverUrl} alt={track.title} />
            <S.TrackCount>x{count}</S.TrackCount>
        </S.TrackPreviewWrapper>
    );
};
