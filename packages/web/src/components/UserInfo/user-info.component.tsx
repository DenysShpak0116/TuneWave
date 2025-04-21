import { FC } from "react";
import { IUser } from "types/user/user.type";
import { UserBlock, Wrapper, Avatar, Name, Stats, Bio, SettingsIcon } from "./user-info.style";
import { FiSettings } from "react-icons/fi";


interface IUserInfoProps {
    user: IUser;
    isMainUser: boolean
}

export const UserInfo: FC<IUserInfoProps> = ({ user, isMainUser }) => {
    return (
        <Wrapper>
            <UserBlock>
                <div style={{ display: "flex", alignItems: "center", gap: "16px" }}>
                    <Avatar src={user.profilePictureUrl} alt={user.username} />
                    <div>
                        <Name>{user.username}</Name>
                        <Stats>
                            0 підписник(ів)/0 підписка(ок)/0 плейлист(ів)
                        </Stats>
                        <Bio>
                            <strong>О собі:</strong> {user.profileInfo || "Юзер забув це заповнити"}
                        </Bio>
                    </div>
                </div>
                {isMainUser && (
                    <SettingsIcon>
                        <FiSettings size={20} />
                    </SettingsIcon>
                )}
            </UserBlock>
        </Wrapper>
    );
};
