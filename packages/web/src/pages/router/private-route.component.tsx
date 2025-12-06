import { useAuthStore } from "@modules/LoginForm/store/store";
import { FC } from "react";
import { Navigate, Outlet } from "react-router-dom";
import { ROUTES } from "./consts/routes.const";

export const PrivateRoute: FC = () => {
    const isAuth = useAuthStore(state => state.isAuth);
    return isAuth() ? <Outlet /> : <Navigate to={ROUTES.SIGN_IN} replace />;
}