import { SongDetailsInputType } from "@modules/CreateTrackForm/const/song-details.const";
import { FC } from "react";
import { InputGroup, InputsWrapper, SongDetailsContainer, StyledInput } from "./song-details.style";
import { MultiInput } from "@ui/MultiInput/multi-input.component";
import { Button } from "@ui/Btn/btn.component";

interface ISongDetailsProps {
    songDetailsInputs: SongDetailsInputType[];
    formData: Record<string, any>;
    setFormData: (data: Record<string, string | File>) => void;
    submitFn: () => void
}


export const SongDetails: FC<ISongDetailsProps> = ({ songDetailsInputs, formData, setFormData, submitFn }) => {
    return (
        <SongDetailsContainer>
            <h3>Деталі треку</h3>
            <InputsWrapper>
                {songDetailsInputs.map((input) => (
                    <InputGroup key={input.name}>
                        <label>{input.placeholder}</label>
                        {input.isMulti ? (
                            <MultiInput
                                name={input.name}
                                placeholder={input.placeholder}
                                value={formData[input.name] || []}
                                onChange={(values) =>
                                    setFormData({ ...formData, [input.name]: values })
                                }
                            />
                        ) : (
                            <StyledInput
                                type={input.type}
                                name={input.name}
                                placeholder={input.placeholder}
                                value={formData[input.name] || ""}
                                onChange={(e) =>
                                    setFormData({ ...formData, [input.name]: e.target.value })
                                }
                            />
                        )}
                    </InputGroup>
                ))}
                <Button
                    onClick={submitFn}
                    text="Загрузити трек"
                    type="button" />
            </InputsWrapper>
        </SongDetailsContainer>
    );
};

