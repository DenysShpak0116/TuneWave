import { FC, ReactNode } from "react";
import { PlayerInsideContainer, StyledPlayerContainer } from "./player-container.component.style";

interface IPlayerContainer {
    children: ReactNode
}

export const PlayerContainer: FC<IPlayerContainer> = ({ children }) => {

    return (
        <StyledPlayerContainer>
            <PlayerInsideContainer>
                {children}
            </PlayerInsideContainer>
        </StyledPlayerContainer>
    );
}