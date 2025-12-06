import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { IChatPreviewType } from "types/chat/chat-preview";
import { ChatPreviewContainer, SearchInput } from "./chat-list.style";
import { ChatPreviewCard } from "@components/ChatPreviewCard/chat-preview-card.component";
import { Button } from "@ui/Btn/btn.component";
import { InviteFollowersModal } from "@modules/InviteFollowersModal";
import { useGetUser } from "pages/UserProfilePage/hooks/useGetUserById";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { useAuthStore } from "@modules/LoginForm/store/store";

interface ChatListProps {
    chatPreviews?: IChatPreviewType[];
    targetUserIds: string[];
}

export const ChatList: FC<ChatListProps> = ({ chatPreviews, targetUserIds }) => {
    const id = useAuthStore(state => state.user?.id);
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null);
    const [isModalOpen, setIsModalOpen] = useState<boolean>(false);
    const { data: user, isLoading } = useGetUser(id!);
    const navigate = useNavigate();

    useEffect(() => {
        if (!chatPreviews || targetUserIds.length === 0) return;

        const match = chatPreviews.find(chat =>
            chat.userIds.some(id => targetUserIds.includes(id))
        );

        if (match) {
            setSelectedChatId(match.id);
        }
    }, [chatPreviews, targetUserIds]);

    const handleChatClick = (chat: IChatPreviewType) => {
        setSelectedChatId(chat.id);
        navigate(`?userIds=${chat.userIds.join(",")}&name=${chat.chatName}`);
    };

    if (isLoading) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <>
            <ChatPreviewContainer>
                <Button
                    text="Створити групу"
                    style={{ marginTop: "10px", padding: "10px" }}
                    onClick={() => setIsModalOpen(true)}
                />

                <SearchInput placeholder="Пошук чату" />

                {chatPreviews?.map((chat) => (
                    <ChatPreviewCard
                        key={chat.id}
                        chat={chat}
                        isSelected={chat.id === selectedChatId}
                        onClick={() => handleChatClick(chat)}
                    />
                ))}
            </ChatPreviewContainer>

            <InviteFollowersModal
                active={isModalOpen}
                setActive={setIsModalOpen}
                followers={user?.followers ?? []}
            />
        </>
    );
};
