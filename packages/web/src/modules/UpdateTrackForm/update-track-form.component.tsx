import { FC, useState, useRef, ChangeEvent } from "react";
import { ISong } from "types/song/song.type";
import {
    FormContainer,
    Container,
    UploadBox,
    UploadContainer,
    Divider,
    UploadIcon,
    PreviewImage,
    HiddenInput,
    Title,
    InputsContainer
} from "./update-track-form.style";
import { FiUpload, FiCheckCircle } from "react-icons/fi";
import { MultiInput } from "@ui/MultiInput/multi-input.component";
import { Button } from "@ui/Btn/btn.component";
import { StyledInput } from "@components/SongDetailsContainer/song-details.style";
import { useUpdateTrack } from "./hooks/useUpdateTrack";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

interface IUpdateTrackFormProps {
    song: ISong;
}

export const UpdateTrackForm: FC<IUpdateTrackFormProps> = ({ song }) => {
    const navigate = useNavigate()
    const [formData, setFormData] = useState<any>({
        title: song.title,
        genre: song.genre,
        artists: song.authors.map(a => a.name),
        tags: song.songTags.map(t => t.name),
        cover: null,
        song: null,
    });

    const [coverPreview, setCoverPreview] = useState<string>(song.coverUrl);
    const [audioFileName, setAudioFileName] = useState<string | null>(null);

    const inputRef = useRef<HTMLInputElement>(null);
    const audioInputRef = useRef<HTMLInputElement>(null);
    const { mutate: updateTrackMutation } = useUpdateTrack();

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

    const handleSubmit = () => {
        const data = new FormData();
        data.append("id", song.id);
        data.append("title", formData.title);
        data.append("genre", formData.genre);
        formData.artists.forEach((artist: string) => data.append("artists", artist));
        formData.tags.forEach((tag: string) => data.append("tags", tag));
        if (formData.song) data.append("song", formData.song);
        if (formData.cover) data.append("cover", formData.cover);

        updateTrackMutation({ songId: song.id, formData: data });
        navigate(ROUTES.TRACK_PAGE.replace(":id", song.id))
    };

    return (
        <FormContainer>
            <Title>Редагування музичного твору</Title>

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

                <UploadContainer>
                    <h3>Файл музичного твору:</h3>
                    <UploadBox onClick={() => audioInputRef.current?.click()}>
                        {audioFileName ? (
                            <div style={{ textAlign: "center" }}>
                                <FiCheckCircle size={36} color="lightgreen" />
                                <span>{audioFileName}</span>
                            </div>
                        ) : (
                            <UploadIcon><FiUpload size={36} /></UploadIcon>
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
            <Divider />
            <InputsContainer>
                <label>Назва треку</label>
                <StyledInput
                    type="text"
                    name="Назва треку"
                    placeholder="Додати назву трека"
                    value={formData.title}
                    onChange={(e) => setFormData({ ...formData, title: e.target.value })}
                />
                <label>Жанр</label>
                <StyledInput
                    type="text"
                    name="Жанр"
                    placeholder="Додати жанр"
                    value={formData.genre}
                    onChange={(e) => setFormData({ ...formData, genre: e.target.value })}
                />
                <label>Виконавці</label>
                <MultiInput
                    placeholder="Виконавці"
                    value={formData.artists}
                    onChange={(values) => setFormData({ ...formData, artists: values })}
                />
                <label>Теги</label>
                <MultiInput
                    placeholder="Теги"
                    value={formData.tags}
                    onChange={(values) => setFormData({ ...formData, tags: values })}
                />
            </InputsContainer>
            <Divider />
            <Button text="Зберегти зміни" onClick={handleSubmit} />
        </FormContainer>
    );
};
