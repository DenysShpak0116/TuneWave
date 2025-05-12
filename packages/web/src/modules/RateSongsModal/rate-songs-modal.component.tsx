import { FC, useState, useMemo } from "react";
import { ISong } from "types/song/song.type";
import { ModalContent, Overlay, SongBlock, SliderBlock, NextButton, ValuesBlock } from "./rate-songs-modal.style";
import { ValueContainer } from "react-select/animated";

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
    const [currentIndex, setCurrentIndex] = useState(0);
    const [sliderValue, setSliderValue] = useState(0);
    const [results, setResults] = useState<IComparisonResult[]>([]);
    const [isAnimating, setIsAnimating] = useState(false);

    const songPairs = useMemo(() => {
        const pairs: { song1: ISong; song2: ISong }[] = [];
        for (let i = 0; i < collectionSongs.length; i++) {
            for (let j = i + 1; j < collectionSongs.length; j++) {
                pairs.push({ song1: collectionSongs[i], song2: collectionSongs[j] });
            }
        }
        return pairs;
    }, [collectionSongs]);

    const isFinished = currentIndex >= songPairs.length;
    const currentPair = songPairs[currentIndex];

    const handleNext = () => {
        setIsAnimating(true);

        setResults((prev) => [
            ...prev,
            {
                song1Id: currentPair.song1.id,
                song2Id: currentPair.song2.id,
                value: sliderValue * -1,
            },
        ]);

        setTimeout(() => {
            setSliderValue(0);
            setCurrentIndex((prev) => prev + 1);
            setIsAnimating(false);
        }, 300);
    };

    const handleFinish = () => {
        console.log("Результати зрівняння:", results);
        setActive(false);
    };

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} $animating={isAnimating} onClick={(e) => e.stopPropagation()}>
                {isFinished ? (
                    <div>
                        <h2>Зрівняння завершено!</h2>
                        <NextButton onClick={handleFinish}>Готово</NextButton>
                    </div>
                ) : (
                    <>
                        <SongBlock>
                            {[currentPair.song1, currentPair.song2].map((song) => (
                                <div key={song.id} className="song-card">
                                    <img src={song.coverUrl} alt={song.title} />
                                    <div className="title">{song.title}</div>
                                    <div className="authors">
                                        {song.authors.map((a) => a.name).join(", ")}
                                    </div>
                                    <div className="genre">Жанр: {song.genre}</div>
                                    <div className="duration">Тривалість: {song.duration}</div>
                                    <div className="stats">
                                        <span>👍 {song.likes}</span>
                                        <span>👎 {song.dislikes}</span>
                                        <span>▶️ {song.listenings}</span>
                                    </div>
                                </div>
                            ))}
                        </SongBlock>

                        <SliderBlock>
                            <ValuesBlock>
                                <span>{currentPair.song1.title} краще</span>
                                <span>{currentPair.song1.title} трохи краще</span>
                                <span>Однаково</span>
                                <span>{currentPair.song2.title} трохи краще</span>
                                <span>{currentPair.song2.title} краще</span>
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