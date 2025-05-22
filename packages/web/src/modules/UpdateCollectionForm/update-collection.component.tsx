import { FC, useState, useRef, ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { FiUpload } from "react-icons/fi";
import { Button } from "@ui/Btn/btn.component";
import { StyledInput } from "@components/SongDetailsContainer/song-details.style";
import { useUpdateTrack } from "./hooks/useUpdateCollection";
import { Container, Divider, FormContainer, HiddenInput, InputsContainer, PreviewImage, Title, UploadBox, UploadContainer, UploadIcon } from "@modules/UpdateTrackForm/update-track-form.style";

interface ICollection {
    id: string;
    title: string;
    description: string;
    coverUrl: string;
}

interface UpdateCollectionFormProps {
    collection: ICollection;
}

export const UpdateCollectionForm: FC<UpdateCollectionFormProps> = ({ collection }) => {
    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        title: collection.title,
        description: collection.description,
        cover: null as File | null,
    });

    const [coverPreview, setCoverPreview] = useState<string>(collection.coverUrl);
    const inputRef = useRef<HTMLInputElement>(null);
    const { mutate: updateCollectionMutation } = useUpdateTrack();

    const handleCoverUpload = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setCoverPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
            setFormData((prev) => ({ ...prev, cover: file }));
        }
    };

    const handleSubmit = () => {
        const data = new FormData();
        data.append("id", collection.id);
        data.append("title", formData.title);
        data.append("description", formData.description);
        if (formData.cover) data.append("cover", formData.cover);

        updateCollectionMutation({ collectionId: collection.id, formData: data });
        navigate(ROUTES.COLLECTION_PAGE.replace(":id", collection.id));
    };

    return (
        <FormContainer>
            <Title>Редагування колекції</Title>

            <Container>
                <UploadContainer>
                    <h3>Обкладинка:</h3>
                    <UploadBox onClick={() => inputRef.current?.click()}>
                        {coverPreview ? (
                            <PreviewImage src={coverPreview} alt="preview" />
                        ) : (
                            <UploadIcon><FiUpload size={36} /></UploadIcon>
                        )}
                        <HiddenInput
                            type="file"
                            accept="image/*"
                            onChange={handleCoverUpload}
                            ref={inputRef}
                        />
                    </UploadBox>
                </UploadContainer>
            </Container>

            <Divider />

            <InputsContainer>
                <label>Назва колекції</label>
                <StyledInput
                    type="text"
                    placeholder="Введіть назву"
                    value={formData.title}
                    onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                />
                <label>Опис</label>
                <StyledInput
                    type="text"
                    placeholder="Введіть опис"
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                />
            </InputsContainer>

            <Divider />

            <Button text="Зберегти зміни" onClick={handleSubmit} />
        </FormContainer>
    );
};
