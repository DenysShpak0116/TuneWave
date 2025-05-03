import {
    BrowserRouter,
    Routes,
    Route,
} from "react-router-dom";

import { ROUTES } from "./consts/routes.const";
import { HomePage } from "pages/HomePage/home.page";
import { LoginPage } from "pages/LoginPage/login.page";
import { RegistrationPage } from "pages/RegistrationPage/registration.page";
import { CreateTrackPage } from "pages/CreateTrackPage/create-track.page";
import { UserProfilePage } from "pages/UserProfilePage/user-profile.page";
import { TrackPage } from "pages/TrackPage/track.page";
import { NotFoundPage } from "pages/NotFoundPage/not-found.page";
import { UpdateUserPage } from "pages/UpdateUserPage/update-user.page";
import { Player } from "@modules/Player";

export const AppRoutes = () => {
    return (
        <BrowserRouter>
            <Player />
            <Routes>
                <Route path={ROUTES.HOME} element={<HomePage />} />
                <Route path={ROUTES.SIGN_IN} element={<LoginPage />} />
                <Route path={ROUTES.SIGN_UP} element={<RegistrationPage />} />
                <Route path={ROUTES.CREATE_TRACK} element={<CreateTrackPage />} />
                <Route path={ROUTES.USER_PROFILE} element={<UserProfilePage />} />
                <Route path={ROUTES.TRACK_PAGE} element={<TrackPage />} />
                <Route path={ROUTES.UPDATE_USER_PAGE} element={<UpdateUserPage />} />
                <Route path="*" element={<NotFoundPage />} />
            </Routes>
        </BrowserRouter>
    );
};
