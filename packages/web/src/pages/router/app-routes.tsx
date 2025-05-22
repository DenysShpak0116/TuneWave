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
import { CollectionPage } from "pages/CollectionPage/collection.page";
import { UpdateTrackPage } from "pages/UpdateTrackPage/update-track.page";
import { AddCriterionPage } from "pages/AddCriterionPage/add-criterion.page";
import { CollectionSongsPage } from "pages/CollectionSongsCriterionsPage/CollectionSongsCritetion.page";
import { CollectiveDecisionPage } from "pages/CollectiveDecisionPage/collective-decision.page";
import { ForgotPasswordPage } from "pages/ForgotPasswordPage/forgot-password.page";
import { ResetPasswordPage } from "pages/ResetPasswordPage/reset-password.page";
import { ChatPage } from "pages/ChatPage/chat.page";
import { GenrePage } from "pages/GenrePage/genre.page";
import { CollectionsPage } from "pages/CollectionsPage/collectionsPage";
import { GenreSongsPage } from "pages/GenreSongsPage/genre-songs.page";
import { UpdateCollectionPage } from "pages/UpdateCollectionPage/update-collection.page";

export const AppRoutes = () => {
    return (
        <BrowserRouter>
            <Player />
            <Routes>
                <Route path={ROUTES.HOME} element={<HomePage />} />
                <Route path={ROUTES.SIGN_IN} element={<LoginPage />} />
                <Route path={ROUTES.SIGN_UP} element={<RegistrationPage />} />
                <Route path={ROUTES.FORGOT_PASSWORD_PAGE} element={<ForgotPasswordPage />} />
                <Route path={ROUTES.CREATE_TRACK} element={<CreateTrackPage />} />
                <Route path={ROUTES.USER_PROFILE} element={<UserProfilePage />} />
                <Route path={ROUTES.TRACK_PAGE} element={<TrackPage />} />
                <Route path={ROUTES.COLLECTION_PAGE} element={<CollectionPage />} />
                <Route path={ROUTES.UPDATE_USER_PAGE} element={<UpdateUserPage />} />
                <Route path={ROUTES.UPDATE_TRACK_PAGE} element={<UpdateTrackPage />} />
                <Route path={ROUTES.ADD_CRITERION_PAGE} element={<AddCriterionPage />} />
                <Route path={ROUTES.SONGS_CRITERIONS_PAGE} element={<CollectionSongsPage />} />
                <Route path={ROUTES.COLLECTIVE_DECISION_PAGE} element={<CollectiveDecisionPage />} />
                <Route path={ROUTES.RESET_PASSWORD_PAGE} element={<ResetPasswordPage />} />
                <Route path={ROUTES.CHAT_PAGE} element={<ChatPage />} />
                <Route path={ROUTES.GENRE_PAGE} element={<GenrePage />} />
                <Route path={ROUTES.COLLECTIONS_PAGE} element={<CollectionsPage />} />
                <Route path={ROUTES.GENRE_SONGS} element={<GenreSongsPage />} />
                <Route path={ROUTES.UPDATE_COLLECTION} element={<UpdateCollectionPage />} />
                <Route path="*" element={<NotFoundPage />} />
            </Routes>
        </BrowserRouter>
    );
};
