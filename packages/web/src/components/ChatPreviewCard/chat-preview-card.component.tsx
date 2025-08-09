import { IChatPreviewType } from "types/chat/chat-preview";
import { Avatar, CardContainer, Content, Message, Username } from "./chat-preview-card.style";
import { FC } from "react";

interface Props {
    chat: IChatPreviewType;
    isSelected?: boolean;
    onClick?: () => void;
}

export const ChatPreviewCard: FC<Props> = ({ chat, isSelected = false, onClick }) => {
    return (
        <CardContainer isSelected={isSelected} onClick={onClick}>
            <Avatar src={chat.userAvatar} alt="avatar" />
            <Content>
                <Username>@{chat.username}</Username>
                <Message>{chat.lastMessage}</Message>
            </Content>
        </CardContainer>
    );
};
  