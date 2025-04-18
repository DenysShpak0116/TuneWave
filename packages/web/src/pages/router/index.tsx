import { createBrowserRouter } from "react-router-dom";
import { ROUTES } from "./consts/routes.const";
import { HomePage } from "pages/HomePage/home.page";
import { LoginPage } from "pages/LoginPage/login.page";
import { RegistrationPage } from "pages/RegistrationPage/registration.page";

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
    }
])

export default router