import { FC } from "react";
import { Logo } from "./track-logo.style";

interface ITrackLogo {
    logo: string | undefined
}
export const TrackLogo: FC<ITrackLogo> = ({ logo }) => {

    return (
        <Logo src={logo} />
    )
} 