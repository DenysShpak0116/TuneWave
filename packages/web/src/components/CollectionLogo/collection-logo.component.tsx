import { FC } from "react";
import { Logo, LogoContainer } from "./collection-logo.style";

interface ICollectionLogoProps {
    logo: string | undefined;
}

export const CollectionLogo: FC<ICollectionLogoProps> = ({ logo }) => {

    return (
        <LogoContainer>
            <Logo src={logo} />
        </LogoContainer>
    );
};
