import { FC, useState } from "react"
import { FollowType } from "types/user/follow.type"
import { Avatar, ButtonGroup, DeleteBtn, EmptyText, TabButton, UserCard, UserListContainer, Username, Wrapper } from "./user-list.style"
import { useNavigate } from "react-router-dom"
import { ROUTES } from "pages/router/consts/routes.const"
import { useUnfollow } from "./hooks/useDeleteUser"

interface IUserListProps {
    followings: FollowType[]
    followers: FollowType[]
    refetchFn: () => void
}

export const UserList: FC<IUserListProps> = ({ followers, followings, refetchFn }) => {
    const navigate = useNavigate()
    const [activeTab, setActiveTab] = useState<"followers" | "followings">("followers")
    const { mutate: unfollow } = useUnfollow()
    const list = activeTab === "followers" ? followers : followings

    const unfollowUser = async (userId: string) => {
        unfollow(userId)
        refetchFn()
    }

    return (
        <Wrapper>
            <ButtonGroup>
                <TabButton
                    active={activeTab === "followers"}
                    onClick={() => setActiveTab("followers")}
                >
                    Підписники ({followers.length})
                </TabButton>
                <TabButton
                    active={activeTab === "followings"}
                    onClick={() => setActiveTab("followings")}
                >
                    Підписки ({followings.length})
                </TabButton>
            </ButtonGroup>

            <UserListContainer>
                {list.length === 0 ? (
                    <EmptyText>Список порожній</EmptyText>
                ) : (
                    list.map((user) => (
                        <UserCard key={user.id} onClick={() => navigate(ROUTES.USER_PROFILE.replace(":id", user.id))}>
                            <Avatar src={user.profilePictureUrl} alt={user.username} />
                            <Username>{user.username}</Username>
                            {activeTab === "followings" && (
                                <DeleteBtn onClick={(e) => { e.stopPropagation(); unfollowUser(user.id) }}>Відписатись</DeleteBtn>
                            )}
                        </UserCard>
                    ))
                )}
            </UserListContainer>
        </Wrapper>
    )
}
