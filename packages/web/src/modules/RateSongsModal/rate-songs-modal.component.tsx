import { FC, useState, useMemo, useEffect } from "react";
import { ISong } from "types/song/song.type";
import { ModalContent, Overlay, SongBlock, SliderBlock, NextButton, ValuesBlock, ComparisonTable, Leaderboard } from "./rate-songs-modal.style";
import { useCreateResult, useGetUserResults } from "./hooks/useResults";
import { IResultType } from "types/results/result.type";
import { ICompareType } from "types/results/compare.type";
import { getComparisonSymbol } from "./helpers/getComprarsionSymbol";
import { useGetVectorBySongId } from "@modules/AddCriterionToSongModal/hooks/useAddVector";

interface IRateSongsModalProps {
    active: boolean;
    setActive: (value: boolean) => void;
    collectionSongs: ISong[];
    collectionId: string;
}

interface IComparisonResult {
    song1Id: string;
    song2Id: string;
    value: number;
}

export const RateSongsModal: FC<IRateSongsModalProps> = ({ active, setActive, collectionId, collectionSongs }) => {
    const [currentIndex, setCurrentIndex] = useState<number>(0);
    const [sliderValue, setSliderValue] = useState<number>(0);
    const [results, setResults] = useState<IComparisonResult[]>([]);
    const [isAnimating, setIsAnimating] = useState<boolean>(false);
    const [isTopOpen, setIsTopOpen] = useState<boolean>(false);
    const createResultMutation = useCreateResult();
    const { data: userResults, isLoading: isResultsLoading, refetch } = useGetUserResults(collectionId);



    const songPairs = useMemo(() => {
        const pairs: { song1: ISong; song2: ISong }[] = [];
        for (let i = 0; i < collectionSongs.length; i++) {
            for (let j = i + 1; j < collectionSongs.length; j++) {
                pairs.push({ song1: collectionSongs[i], song2: collectionSongs[j] });
            }
        }
        return pairs;
    }, [collectionSongs]);


    const comparisonMatrix = useMemo(() => {
        const matrix: string[][] = Array(collectionSongs.length).fill(null).map(() =>
            Array(collectionSongs.length).fill("=")
        );

        results.forEach((result) => {
            const index1 = collectionSongs.findIndex(s => s.id === result.song1Id);
            const index2 = collectionSongs.findIndex(s => s.id === result.song2Id);

            const symbol = getComparisonSymbol(result.value);
            const invertedSymbol = getComparisonSymbol(result.value * -1);

            if (index1 !== -1 && index2 !== -1) {
                matrix[index1][index2] = symbol;
                matrix[index2][index1] = invertedSymbol;
            }
        });

        return matrix;
    }, [results, collectionSongs]);

    const isFinished = currentIndex >= songPairs.length;
    const currentPair = !isFinished ? songPairs[currentIndex] : null;

    const handleNext = () => {
        setIsAnimating(true);

        setResults((prev) => [
            ...prev,
            {
                song1Id: currentPair!.song1.id,
                song2Id: currentPair!.song2.id,
                value: sliderValue * -1,
            },
        ]);

        setTimeout(() => {
            setSliderValue(0);
            setCurrentIndex((prev) => prev + 1);
            setIsAnimating(false);
        }, 300);
    };

    useEffect(() => {
        if (isFinished) {
            handleFinish();
        }
    }, [isFinished]);


    const handleFinish = () => {
        const preparedResults: IResultType[] = collectionSongs.map((song) => {
            const comparedTo: ICompareType[] = results
                .filter((r) => r.song1Id === song.id || r.song2Id === song.id)
                .map((r) => {
                    if (r.song1Id === song.id) {
                        return { song2Id: r.song2Id, result: r.value };
                    } else {
                        return { song2Id: r.song1Id, result: -r.value };
                    }
                });

            return { song1Id: song.id, comparedTo };
        });

        createResultMutation.mutate({ collectionId, results: preparedResults });
        refetch()
    }



    const { data: vectorsSong1, isLoading: isVectorsSong1Loading } = useGetVectorBySongId(
        collectionId,
        currentPair?.song1.id || ""
    );
    const { data: vectorsSong2, isLoading: isVectorsSong2Loading } = useGetVectorBySongId(
        collectionId,
        currentPair?.song2.id || ""
    );


    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} $animating={isAnimating} onClick={(e) => e.stopPropagation()}>
                {isFinished ? (
                    <div>
                        <h2>Зрівняння завершено!</h2>

                        <ComparisonTable>
                            <thead>
                                <tr>
                                    <th></th>
                                    {collectionSongs.map((song, idx) => (
                                        <th key={idx}>{song.title}</th>
                                    ))}
                                </tr>
                            </thead>
                            <tbody>
                                {comparisonMatrix.map((row, rowIndex) => (
                                    <tr key={rowIndex}>
                                        <th>{collectionSongs[rowIndex].title}</th>
                                        {row.map((cell, colIndex) => (
                                            <td key={colIndex}>{cell}</td>
                                        ))}
                                    </tr>
                                ))}
                            </tbody>
                        </ComparisonTable>

                        <NextButton onClick={() => setIsTopOpen((prev) => !prev)}>
                            {isTopOpen ? "Сховати топ пісень" : "Передивитись топ пісень"}
                        </NextButton>

                        <NextButton onClick={() => setActive(false)}>Готово</NextButton>

                        {isTopOpen && (
                            <div>
                                <h3>🏆 Топ пісень</h3>
                                {isResultsLoading ? (
                                    <div>Завантаження...</div>
                                ) : (
                                    <Leaderboard>
                                        {userResults
                                            .sort((a, b) => a.songRang - b.songRang)
                                            .map((song, index) => (
                                                <li key={song.songId}>
                                                    <span className="medal">
                                                        {index === 0 ? "🥇" : index === 1 ? "🥈" : index === 2 ? "🥉" : `${index + 1}.`}
                                                    </span>
                                                    {song.songName}
                                                </li>
                                            ))}
                                    </Leaderboard>
                                )}
                            </div>
                        )}
                    </div>
                ) : (
                    <>
                        <SongBlock>
                            {[currentPair?.song1, currentPair?.song2].map((song, idx) => {
                                if (!song) return null;

                                const vectors = idx === 0 ? vectorsSong1 : vectorsSong2;
                                const isLoading = idx === 0 ? isVectorsSong1Loading : isVectorsSong2Loading;

                                return (
                                    <div key={song.id} className="song-card">
                                        <img src={song.coverUrl} alt={song.title} />
                                        <div className="title">{song.title}</div>
                                        <div className="authors">
                                            {song.authors.map((a) => a.name).join(", ")}
                                        </div>
                                        <div className="genre">Жанр: {song.genre}</div>
                                        <div className="duration">Тривалість: {song.duration}</div>
                                        <div className="criteria">
                                            {isLoading ? (
                                                <div>Завантаження критеріїв...</div>
                                            ) : (
                                                vectors?.map((v) => (
                                                    <p key={v.criterionId}>
                                                        {v.criterion}: {v.mark}
                                                    </p>
                                                ))
                                            )}
                                        </div>
                                        <div className="stats">
                                            <span>👍 {song.likes}</span>
                                            <span>👎 {song.dislikes}</span>
                                            <span>▶️ {song.listenings}</span>
                                        </div>
                                    </div>
                                );
                            })}
                        </SongBlock>
                        <SliderBlock>
                            <ValuesBlock>
                                <span>{currentPair!.song1.title} краще</span>
                                <span>{currentPair!.song1.title} трохи краще</span>
                                <span>Однаково</span>
                                <span>{currentPair!.song2.title} трохи краще</span>
                                <span>{currentPair!.song2.title} краще</span>
                            </ValuesBlock>
                            <input
                                type="range"
                                min={-2}
                                max={2}
                                step={1}
                                value={sliderValue}
                                onChange={(e) => setSliderValue(Number(e.target.value))}
                            />
                        </SliderBlock>

                        <NextButton onClick={handleNext}>Далі</NextButton>
                    </>
                )}
            </ModalContent>
        </Overlay>
    );
};