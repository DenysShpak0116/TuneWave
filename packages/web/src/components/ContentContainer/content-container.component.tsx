import { ChangeEvent, Dispatch, FC, useRef, useState } from "react";
import {
    Container,
    UploadBox,
    Divider,
    UploadIcon,
    UploadContainer,
    PreviewImage,
    HiddenInput,
} from "./content-container.style";
import { FiCheckCircle, FiUpload } from "react-icons/fi";

interface Props {
    setFormData: Dispatch<React.SetStateAction<Record<string, string | File>>>;
}

export const ContentContainer: FC<Props> = ({ setFormData }) => {
    const [coverPreview, setCoverPreview] = useState<string | null>(null);
    const [audioFileName, setAudioFileName] = useState<string | null>(null);
    const inputRef = useRef<HTMLInputElement>(null);
    const audioInputRef = useRef<HTMLInputElement>(null);

    const handleCoverUpload = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setCoverPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
            setFormData(prev => ({ ...prev, cover: file }));
        }
    };


    const handleAudioUpload = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setAudioFileName(file.name);
            setFormData(prev => ({ ...prev, song: file }));
        }
    };

    return (
        <Container>
            <UploadContainer>
                <h3>Завантажити обкладинку</h3>
                <UploadBox onClick={() => inputRef.current?.click()}>
                    {coverPreview ? (
                        <PreviewImage src={coverPreview} alt="preview" />
                    ) : (
                        <UploadIcon>
                            <FiUpload size={36} />
                        </UploadIcon>
                    )}
                    <HiddenInput
                        type="file"
                        accept="image/*"
                        onChange={handleCoverUpload}
                        ref={inputRef}
                    />
                </UploadBox>
            </UploadContainer>

            <Divider />

            <UploadContainer>
                <h3>Завантажити файл пісні</h3>
                <UploadBox onClick={() => audioInputRef.current?.click()}>
                    {audioFileName ? (
                        <div style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
                            <FiCheckCircle size={36} color="lightgreen" />
                            <span style={{ marginTop: "8px", fontSize: "14px" }}>{audioFileName}</span>
                        </div>
                    ) : (
                        <UploadIcon>
                            <FiUpload size={36} />
                        </UploadIcon>
                    )}
                    <HiddenInput
                        type="file"
                        accept="audio/*"
                        onChange={handleAudioUpload}
                        ref={audioInputRef}
                    />
                </UploadBox>
            </UploadContainer>
        </Container>
    );
};
