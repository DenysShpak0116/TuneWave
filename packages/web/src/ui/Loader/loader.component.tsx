import { FC } from "react";
import { Overlay, Spinner } from "./loader.style";

export const Loader: FC = () => {
    return (
        <Overlay>
            <Spinner />
            <p>Завантаження...</p>
        </Overlay>
    );
};