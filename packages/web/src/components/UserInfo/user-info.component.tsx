import { FC } from "react";
import { IUser } from "types/user/user.type";
import { UserBlock, Wrapper, Avatar, Name, Stats, Bio, SettingsIcon } from "./user-info.style";
import { FiSettings } from "react-icons/fi";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { Button } from "@ui/Btn/btn.component";
import { useFollow, useIsFollowed } from "./hooks/useFollowing";
import { useUnfollow } from "@modules/UserList/hooks/useDeleteUser";

interface IUserInfoProps {
    user: IUser;
    isMainUser: boolean;
    collectionsCount: number;
}

export const UserInfo: FC<IUserInfoProps> = ({ user, isMainUser, collectionsCount }) => {
    const navigate = useNavigate()
    const { data: isFollowed, isLoading, refetch } = useIsFollowed(user.id)
    const { mutate: follow } = useFollow(() => refetch())
    const { mutate: unfollow } = useUnfollow(() => refetch())

    const handleFollow = (id: string) => {
        follow(id)
    }

    const handleUnfollow = (id: string) => {
        unfollow(id);
    }

    return (
        <Wrapper>
            <UserBlock>
                <div style={{ display: "flex", alignItems: "center", gap: "16px" }}>
                    <Avatar src={user.profilePictureUrl} alt={user.username} />
                    <div>
                        <Name>{user.username}</Name>
                        <Stats>
                            {user.followers.length} підписників/ {user.follows.length} підписок/ {collectionsCount} плейлистів
                        </Stats>
                        <Bio>
                            <strong>О собі:</strong> {user.profileInfo || "Юзер забув це заповнити"}
                        </Bio>
                        {!isMainUser && (
                            <>
                                <Button
                                    text="Надіслати повідомлення"
                                    style={{ padding: "5px", marginTop: "10px" }}
                                    onClick={() => navigate(`${ROUTES.CHAT_PAGE}?targetUserId=${user.id}`)}
                                />
                                {!isLoading && !isFollowed && (
                                    <Button
                                        text="Підписатися"
                                        style={{ padding: "5px", marginTop: "10px" }}
                                        onClick={() => {
                                            handleFollow(user.id)
                                        }}
                                    />
                                )}
                                {isFollowed && (
                                    <Button
                                        text="Відписатись"
                                        style={{ padding: "5px", marginTop: "10px", backgroundColor: "rgb(236, 94, 120)" }}
                                        onClick={() => {
                                            handleUnfollow(user.id)
                                        }}
                                    />
                                )}
                            </>
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
