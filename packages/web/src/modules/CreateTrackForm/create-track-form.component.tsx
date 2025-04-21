import { FC, useState } from "react";
import { FormContainer } from "./create-track-form.style";
import { ContentContainer } from "@components/ContentContainer/content-container.component";
import { SongDetails } from "@components/SongDetailsContainer/song-details.component";
import { songDetailsInputs } from "./const/song-details.const";
import { useCreateTrack } from "./hooks/useCreateTrack";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { useAuthStore } from "@modules/LoginForm/store/store";

export const CreateTrackForm: FC = () => {
    const userId = useAuthStore.getState().user?.id;
    const navigate = useNavigate()
    const [formData, setFormData] = useState<Record<string, string | File>>({});
    const { mutate: createTrack } = useCreateTrack();

    const handleSubmit = () => {
        const userId = useAuthStore.getState().user?.id;

        const data = new FormData();

        for (const key in formData) {
            if (formData[key]) {
                data.append(key, formData[key]);
            }
        }
        
        if (userId) {
            data.append("userID", userId);
        }

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
