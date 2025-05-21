import { FC } from "react";
import { IUser } from "types/user/user.type";
import { UserBlock, Wrapper, Avatar, Name, Stats, Bio, SettingsIcon } from "./user-info.style";
import { FiSettings } from "react-icons/fi";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { Button } from "@ui/Btn/btn.component";


interface IUserInfoProps {
    user: IUser;
    isMainUser: boolean
}

export const UserInfo: FC<IUserInfoProps> = ({ user, isMainUser }) => {
    const navigate = useNavigate()

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
                        {!isMainUser && (
                            <Button
                                text="Надіслати повідомлення"
                                style={{ padding: "5px", marginTop: "10px" }}
                                onClick={() => navigate(`${ROUTES.CHAT_PAGE}?targetUserId=${user.id}`)}/>
                        )}

                    </div>
                </div>
                {isMainUser && (
                    <SettingsIcon>
                        <FiSettings onClick={() => navigate(ROUTES.UPDATE_USER_PAGE.replace(':id', user.id))} size={20} />
                    </SettingsIcon>
                )}
            </UserBlock>
        </Wrapper>
    );
};
