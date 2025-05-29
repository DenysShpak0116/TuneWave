import { FC } from "react";
import { MainLayout } from "@ui/layout/main-layout";
import { useNavigate } from "react-router-dom";
import { Button } from "@ui/Btn/btn.component";
import { ROUTES } from "pages/router/consts/routes.const";
import { Subtitle, Title, Wrapper } from "./not-found.style";

export const NotFoundPage: FC = () => {
    const navigate = useNavigate();

    return (
        <MainLayout>
            <Wrapper>
                <Title>404</Title>
                <Subtitle>Сторінку не знайдено</Subtitle>
                <Button
                    text="На головну"
                    onClick={() => navigate(ROUTES.HOME)}
                />
            </Wrapper>
        </MainLayout>
    );
};

