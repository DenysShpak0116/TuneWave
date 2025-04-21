import { createBrowserRouter } from "react-router-dom";
import { ROUTES } from "./consts/routes.const";
import { HomePage } from "pages/HomePage/home.page";
import { LoginPage } from "pages/LoginPage/login.page";
import { RegistrationPage } from "pages/RegistrationPage/registration.page";
import { CreateTrackPage } from "pages/CreateTrackPage/create-track.page";

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
    }
])

export default router