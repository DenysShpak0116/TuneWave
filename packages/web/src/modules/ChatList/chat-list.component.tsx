import { FC, useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { IChatPreviewType } from "types/chat/chat-preview";
import { ChatPreviewContainer, SearchInput } from "./chat-list.style";
import { ChatPreviewCard } from "@components/ChatPreviewCard/chat-preview-card.component";

interface ChatListProps {
    chatPreviews?: IChatPreviewType[];
    targetUserId: string;
}

export const ChatList: FC<ChatListProps> = ({ chatPreviews, targetUserId }) => {
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        if (!chatPreviews) return;

        const match = chatPreviews.find(chat =>
            chat.targetUserId === targetUserId
        );

        if (match) {
            setSelectedChatId(match.id);
        }
    }, [chatPreviews, targetUserId]);

    const handleChatClick = (chat: IChatPreviewType) => {
        setSelectedChatId(chat.id);
        navigate(`?targetUserId=${chat.targetUserId}`);
    };

    return (
        <ChatPreviewContainer>
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
    );
};
