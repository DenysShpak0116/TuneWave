import { FC } from "react";
import { LogoutBlock, UserBlockContainer, UserProfileImg } from "./user-block.style";
import logoutIcon from "@assets/images/ic_logout.png"

interface IUserBlockProps {
    profileImg: string | undefined;
    username: string | undefined;
    logoutFn: VoidFunction;
    userProfileFn: VoidFunction
}

export const UserBlock: FC<IUserBlockProps> = ({ profileImg, username, logoutFn, userProfileFn }) => {

    return (
        <UserBlockContainer>

            <span>{username}</span>
            <UserProfileImg onClick={userProfileFn} src={profileImg} alt="user-image" />

            <LogoutBlock src={logoutIcon} onClick={logoutFn} />
        </UserBlockContainer>
    )
}