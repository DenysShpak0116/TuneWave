import { FC, useMemo } from "react";
import { useGetListensByDay, useGetAvarageListens, useGetMedianListens, useGetPeakAndSilentDay } from "./hooks/useGetStatistic";
import {
    ResponsiveContainer,
    LineChart,
    Line,
    CartesianGrid,
    XAxis,
    YAxis,
    Tooltip,
    Legend,
    BarChart,
    Bar,
    PieChart,
    Pie,
    Cell,
} from "recharts";
import * as S from "./statistic-block.style";
import { Loader } from "@ui/Loader/loader.component";

export const StatisticBlock: FC = () => {
    const { data, isLoading, isError } = useGetListensByDay();
    const { data: avgData } = useGetAvarageListens();
    const { data: medianData } = useGetMedianListens();
    const { data: peakData } = useGetPeakAndSilentDay();

    const chartData = useMemo(() => {
        if (!data) return [];
        return Object.entries(data)
            .map(([day, listens]) => ({
                day,
                listens: Number(listens),
            }))
            .sort((a, b) => new Date(a.day).getTime() - new Date(b.day).getTime());
    }, [data]);

    if (isLoading) return <Loader />;
    if (isError) return <S.Message>Помилка при завантаженні даних</S.Message>;

    const totalListens = chartData.reduce((acc, d) => acc + d.listens, 0);

    return (
        <S.Container>
            <S.Title>📊 Статистика прослуховування</S.Title>

            <S.StatsRow>
                <S.StatCard>
                    <S.StatValue>{avgData?.avg?.toFixed(2) ?? "-"}</S.StatValue>
                    <S.StatLabel>Середнє значення</S.StatLabel>
                </S.StatCard>

                <S.StatCard>
                    <S.StatValue>{medianData?.median?.toFixed(2) ?? "-"}</S.StatValue>
                    <S.StatLabel>Медіана</S.StatLabel>
                </S.StatCard>

                <S.StatCard>
                    <S.StatValue>{peakData?.peak_count ?? "-"}</S.StatValue>
                    <S.StatLabel>Пік активності ({peakData?.peak_day})</S.StatLabel>
                </S.StatCard>

                <S.StatCard>
                    <S.StatValue>{peakData?.silent_count ?? "-"}</S.StatValue>
                    <S.StatLabel>Найспокійніший день ({peakData?.silent_day})</S.StatLabel>
                </S.StatCard>
            </S.StatsRow>

            <S.ChartWrapper>
                <S.ChartCard>
                    <h3>Динаміка прослуховувань</h3>
                    <ResponsiveContainer width="100%" height={300}>
                        <LineChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="day" />
                            <YAxis />
                            <Tooltip />
                            <Legend />
                            <Line
                                type="monotone"
                                dataKey="listens"
                                stroke="#4A90E2"
                                strokeWidth={3}
                                dot={{ r: 4 }}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                </S.ChartCard>

                <S.ChartCard>
                    <h3>Порівняння по дням</h3>
                    <ResponsiveContainer width="100%" height={300}>
                        <BarChart data={chartData}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="day" />
                            <YAxis />
                            <Tooltip />
                            <Bar
                                dataKey="listens"
                                fill="#82ca9d"
                                barSize={40}
                                radius={[8, 8, 0, 0]}
                            />
                        </BarChart>
                    </ResponsiveContainer>
                </S.ChartCard>

                <S.ChartCard>
                    <h3>Відсотковий розподіл</h3>
                    <ResponsiveContainer width="100%" height={300}>
                        <PieChart>
                            <Pie
                                data={chartData}
                                dataKey="listens"
                                nameKey="day"
                                outerRadius={110}
                                label={({ percent }) => `${(percent! * 100).toFixed(1)}%`}
                            >
                                {chartData.map((_, i) => (
                                    <Cell
                                        key={i}
                                        fill={["#4A90E2", "#50E3C2", "#B8E986", "#F5A623", "#D0021B"][i % 5]}
                                    />
                                ))}
                            </Pie>

                            <Tooltip
                                formatter={(value: number, _name: string) => {
                                    const percent = ((value / totalListens) * 100).toFixed(1);
                                    return [`${value} (${percent}%)`, `Прослуховувань за дату: ${_name}`];
                                }}
                                labelFormatter={(day: string) => `Дата: ${day}`}
                            />
                        </PieChart>
                    </ResponsiveContainer>

                    <S.Total>Всього прослуховувань: {totalListens}</S.Total>
                </S.ChartCard>
            </S.ChartWrapper>
        </S.Container>
    );
};
