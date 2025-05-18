import { FC, useState } from "react";
import { IChatPreviewType } from "types/chat/chat-preview";
import { ChatPreviewContainer, SearchInput } from "./chat-list.style";
import { ChatPreviewCard } from "@components/ChatPreviewCard/chat-preview-card.component";

interface ChatListProps {
    chatPreviews?: IChatPreviewType[];
}

export const ChatList: FC<ChatListProps> = ({ chatPreviews = [] }) => {
    const [selectedChatId, setSelectedChatId] = useState<string | null>(null);

    return (
        <ChatPreviewContainer>
            <SearchInput placeholder="Пошук чату" />
            {chatPreviews.map((chat) => (
                <ChatPreviewCard
                    key={chat.id}
                    chat={chat}
                    isSelected={chat.id === selectedChatId}
                    onClick={() => setSelectedChatId(chat.id)}
                />
            ))}
        </ChatPreviewContainer>
    );
};