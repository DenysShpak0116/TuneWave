import { FC, ReactNode } from "react";
import { TrackInformationContainer } from "./track-information-layout.style";

interface TrackInformationContainerProps {
    children: ReactNode
}

export const TrackInformationLayout: FC<TrackInformationContainerProps> = ({ children }) => {
    return (
        <TrackInformationContainer>
            {children}
        </TrackInformationContainer>
    )
}