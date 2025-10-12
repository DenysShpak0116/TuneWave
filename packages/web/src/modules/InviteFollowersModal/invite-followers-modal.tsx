import { ModalContent, ModalHeader, ModalHeaderText, Overlay } from "@modules/SelectCollectionModal/select-collection.style";
import { FC, useState } from "react";
import { FollowType } from "types/user/follow.type";
import { Avatar, CreateButton, FollowerItem, ListContainer, Username, ChatNameInput } from "./invite-followers-modal.style";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

interface InviteFollowersModalProps {
    active: boolean;
    setActive: (value: boolean) => void;
    followers: FollowType[];
}

export const InviteFollowersModal: FC<InviteFollowersModalProps> = ({
    active,
    setActive,
    followers,
}) => {
    const [selectedIds, setSelectedIds] = useState<string[]>([]);
    const [chatName, setChatName] = useState<string>("");
    const navigate = useNavigate()

    const toggleSelect = (id: string) => {
        setSelectedIds(prev =>
            prev.includes(id)
                ? prev.filter(i => i !== id)
                : [...prev, id]
        );
    };

    const handleCreate = () => {
        navigate(`${ROUTES.CHAT_PAGE}?userIds=${selectedIds.join(',')}&name=${chatName}`)
    };

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} onClick={(e) => e.stopPropagation()}>
                <ModalHeader>
                    <ModalHeaderText>Оберіть друзів, щоб створити чат</ModalHeaderText>
                </ModalHeader>
                <ChatNameInput
                    placeholder="Назва чату"
                    value={chatName}
                    onChange={(e) => setChatName(e.target.value)}
                />

                <ListContainer>
                    {followers.map((follower) => {
                        const isSelected = selectedIds.includes(follower.id);

                        return (
                            <FollowerItem
                                key={follower.id}
                                $selected={isSelected}
                                onClick={() => toggleSelect(follower.id)}
                            >
                                <Avatar src={follower.profilePictureUrl} />
                                <Username>{follower.username}</Username>
                            </FollowerItem>
                        );
                    })}
                </ListContainer>

                <CreateButton
                    disabled={selectedIds.length === 0 || chatName.trim() === ""}
                    onClick={handleCreate}
                >
                    Створити
                </CreateButton>
            </ModalContent>
        </Overlay>
    );
};
