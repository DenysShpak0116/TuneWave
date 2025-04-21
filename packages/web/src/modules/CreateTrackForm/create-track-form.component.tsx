import { FC, useState } from "react";
import { FormContainer } from "./create-track-form.style";
import { ContentContainer } from "@components/ContentContainer/content-container.component";
import { SongDetails } from "@components/SongDetailsContainer/song-details.component";
import { songDetailsInputs } from "./const/song-details.const";
import { useCreateTrack } from "./hooks/useCreateTrack";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { CreateTrackRequest } from "./types/createTrackRequest.type";
import toast from "react-hot-toast";
import { validateCreateTrackFormData } from "./helpers/trackFormDataValidation";

export const CreateTrackForm: FC = () => {
    const navigate = useNavigate()
    const [formData, setFormData] = useState<Partial<CreateTrackRequest>>({});
    const { mutate: createTrack } = useCreateTrack();

    const handleSubmit = () => {
        const userId = useAuthStore.getState().user?.id;

        if (!userId) {
            toast.error("Користувач не авторизований.");
            return;
        }

        const fullFormData = {
            ...formData,
            userID: userId,
        };

        if (!validateCreateTrackFormData(fullFormData)) {
            return;
        }

        const data = new FormData();
        data.append("userID", fullFormData.userID);
        data.append("title", fullFormData.title);
        data.append("genre", fullFormData.genre);
        fullFormData.artists.forEach(artist => data.append("artists", artist));
        fullFormData.tags.forEach(tag => data.append("tags", tag));
        data.append("song", fullFormData.song);
        data.append("cover", fullFormData.cover);

        createTrack(data);
        navigate(ROUTES.HOME);
    };

    return (
        <FormContainer>
            <ContentContainer
                setFormData={setFormData}
            />
            <SongDetails
                songDetailsInputs={songDetailsInputs}
                formData={formData}
                setFormData={setFormData}
                submitFn={handleSubmit}
            />
        </FormContainer>
    );
};
