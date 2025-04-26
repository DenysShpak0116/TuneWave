import { createBrowserRouter } from "react-router-dom";
import { ROUTES } from "./consts/routes.const";
import { HomePage } from "pages/HomePage/home.page";
import { LoginPage } from "pages/LoginPage/login.page";
import { RegistrationPage } from "pages/RegistrationPage/registration.page";
import { CreateTrackPage } from "pages/CreateTrackPage/create-track.page";
import { UserProfilePage } from "pages/UserProfilePage/user-profile.page";
import { TrackPage } from "pages/TrackPage/track.page";
import { NotFoundPage } from "pages/NotFoundPage/not-found.page";
import { UpdateUserPage } from "pages/UpdateUserPage/update-user.page";

const router = createBrowserRouter([
    {
        path: ROUTES.HOME,
        element: <HomePage />
    },
    {
        path: ROUTES.SIGN_IN,
        element: <LoginPage />
    },
    {
        path: ROUTES.SIGN_UP,
        element: <RegistrationPage />
    },
    {
        path: ROUTES.CREATE_TRACK,
        element: <CreateTrackPage />
    },
    {
        path: ROUTES.USER_PROFILE,
        element: <UserProfilePage />
    },
    {
        path: ROUTES.TRACK_PAGE,
        element: <TrackPage />
    },
    {
        path: ROUTES.UPDATE_USER_PAGE,
        element: <UpdateUserPage />
    },
    {
        path: "*",
        element: <NotFoundPage />
    }
])

export default router