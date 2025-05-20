import { FC, useEffect, useRef, useState } from "react";
import { IMessageType } from "types/chat/message.type";
import { Avatar, Container, Header, Input, InputWrapper, MessageBubble, MessageRow, SendButton, Timestamp, Username, Wrapper } from "./main-chat.style";
import sendIcon from "@assets/images/ic_send.png"

interface MainChatProps {
    messages: IMessageType[];
    currentUserId: string;
    partnerUsername: string;
    partnerAvatar: string;
    onSendMessage: (msg: string) => void
}

export const MainChat: FC<MainChatProps> = ({
    messages,
    currentUserId,
    partnerUsername,
    partnerAvatar,
    onSendMessage
}) => {
    const [input, setInput] = useState("");
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (containerRef.current) {
            containerRef.current.scrollTop = containerRef.current.scrollHeight;
        }
    }, [messages]);

    const handleSend = (e: any) => {
        e.preventDefault();
        if (input.trim()) {
            onSendMessage?.(input.trim());
            setInput("");
        }
    };

    return (
        <Wrapper>
            <Header>
                <Avatar src={partnerAvatar} />
                <Username>@{partnerUsername}</Username>
            </Header>
            <Container ref={containerRef}>
                {messages.map((msg) => {
                    const isCurrentUser = msg.senderId === currentUserId;
                    return (
                        <MessageRow key={msg.id} isCurrentUser={isCurrentUser}>
                            <MessageBubble isCurrentUser={isCurrentUser}>
                                {msg.content}
                                <Timestamp>
                                    {msg.createdAt.slice(11, 16)}
                                </Timestamp>
                            </MessageBubble>
                        </MessageRow>
                    );
                })}
            </Container>
            <form onSubmit={handleSend}>
                <InputWrapper>
                    <Input
                        placeholder="Написати повідомлення..."
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                    />
                    <SendButton type="submit">
                        <img src={sendIcon} />
                    </SendButton>
                </InputWrapper>
            </form>
        </Wrapper>
    );
};