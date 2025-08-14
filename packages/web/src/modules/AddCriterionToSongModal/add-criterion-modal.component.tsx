import { FC, useEffect, useState } from "react";
import { InputsBlock, ModalContent, NextButton, Overlay, SongHeader } from "./add-criterion-modal.style";
import { ISong } from "types/song/song.type";
import { useCriterions } from "pages/AddCriterionPage/hooks/useCriterions";
import { useAddVector, useGetVectorBySongId, useUpdateVector } from "./hooks/useAddVector";
import { IVectorType } from "types/vectors/vector.type";

interface IAddCriterionModalProps {
    active: boolean;
    setActive: (value: boolean) => void;
    collectionSongs: ISong[];
    collectionId: string;
}

export const AddCriterionModal: FC<IAddCriterionModalProps> = ({
    active,
    setActive,
    collectionSongs,
    collectionId,
}) => {
    const { data: criterions } = useCriterions();
    const [currentIndex, setCurrentIndex] = useState(0);
    const [inputs, setInputs] = useState<{ [criterionId: string]: string }>({});
    const [isAnimating, setIsAnimating] = useState(false);

    const updateVectorMutation = useUpdateVector();
    const addVectorMutation = useAddVector();

    const currentSong = collectionSongs[currentIndex];
    const isFinished = currentIndex >= collectionSongs.length;

    const { data: vectorsData } = useGetVectorBySongId(collectionId, currentSong?.id);

    useEffect(() => {
        if (vectorsData && criterions) {
            const initialInputs: { [criterionId: string]: string } = {};
            criterions.forEach((criterion) => {
                const vector = vectorsData.find((v: IVectorType) => v.criterionId === criterion.id);
                if (vector) {
                    initialInputs[criterion.id] = vector.mark;
                }
            });
            setInputs(initialInputs);
        } else {
            setInputs({});
        }
    }, [vectorsData, criterions, currentSong?.id]);

    const allFilled = criterions?.every((criterion) => {
        const input = inputs[criterion.id];
        return typeof input === "string" && input.trim() !== "";
    });

    const handleInputChange = (criterionId: string, value: string) => {
        setInputs((prev) => ({ ...prev, [criterionId]: value }));
    };

    
    const handleNext = () => {
        setIsAnimating(true);

        setTimeout(() => {
            const vectors = Object.entries(inputs).map(([criterionId, value]) => ({
                criterionId,
                mark: value,
            }));

            const existingVectorsMap = new Map(
                (vectorsData || []).map((v) => [v.criterionId, v])
            );

            const toUpdate = vectors
                .filter(v => existingVectorsMap.has(v.criterionId))
                .map(v => {
                    const existingVector = existingVectorsMap.get(v.criterionId) as IVectorType;
                    return {
                        id: existingVector.id,
                        criterionId: v.criterionId,
                        mark: v.mark,
                    };
                });

            const toAdd = vectors.filter(v => !existingVectorsMap.has(v.criterionId));

            let finishedRequests = 0;
            const totalRequests = (toUpdate.length > 0 ? 1 : 0) + (toAdd.length > 0 ? 1 : 0);

            const handleSuccess = () => {
                finishedRequests++;
                if (finishedRequests === totalRequests) {
                    setInputs({});
                    setCurrentIndex((prev) => prev + 1);
                    setIsAnimating(false);
                }
            };

            const handleError = (error: any) => {
                console.error("Ошибка:", error);
                setIsAnimating(false);
            };

            if (toUpdate.length > 0) {
                updateVectorMutation.mutate(
                    {
                        collectionId,
                        songId: currentSong.id,
                        vectors: toUpdate,
                    },
                    {
                        onSuccess: handleSuccess,
                        onError: handleError,
                    }
                );
            }

            if (toAdd.length > 0) {
                addVectorMutation.mutate(
                    {
                        collectionId,
                        songId: currentSong.id,
                        vectors: toAdd,
                    },
                    {
                        onSuccess: handleSuccess,
                        onError: handleError,
                    }
                );
            }

            if (toUpdate.length === 0 && toAdd.length === 0) {
                setIsAnimating(false);
            }

        }, 300);
    };

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} $animating={isAnimating} onClick={(e) => e.stopPropagation()}>
                {isFinished ? (
                    <h2>Усі пісні пройдені!</h2>
                ) : (
                    <>
                        <SongHeader>
                            <p>{currentIndex + 1}/{collectionSongs.length}</p>
                            <img src={currentSong.coverUrl} alt="cover" />
                            <h3>{currentSong.title}</h3>
                        </SongHeader>

                        <InputsBlock>
                            {criterions?.map((criterion) => (
                                <div key={criterion.id}>
                                    <label>{criterion.name}</label>
                                    <input
                                        type="text"
                                        value={inputs[criterion.id] || ""}
                                        onChange={(e) =>
                                            handleInputChange(criterion.id, e.target.value)
                                        }
                                    />
                                </div>
                            ))}
                        </InputsBlock>

                        <NextButton disabled={!allFilled || isAnimating} onClick={handleNext}>
                            Далі
                        </NextButton>
                    </>
                )}
            </ModalContent>
        </Overlay>
    );
};
