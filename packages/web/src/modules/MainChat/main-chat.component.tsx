import { FC, useEffect, useMemo, useRef, useState } from "react";
import CryptoJS from "crypto-js";
import { aesDecrypt } from "@modules/MainChat/hooks/cryptoAes";
import { IMessageType } from "types/chat/message.type";
import {
    Avatar, AvatarChat, Container, Content, Header, Input,
    InputWrapper, MessageBubble, MessageRow, SendButton,
    Timestamp, Username, UsernameChat, Wrapper
} from "./main-chat.style";
import sendIcon from "@assets/images/ic_send.png";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

interface MainChatProps {
    messages: IMessageType[];
    currentUserId: string;
    users: string[];
    chatName: string;
    chatPhoto?: string;
    chatId?: string;
    onSendMessage: (msg: string) => void;
}

export const MainChat: FC<MainChatProps> = ({
    messages,
    currentUserId,
    chatPhoto,
    users,
    chatName,
    chatId,
    onSendMessage,
}) => {
    const navigate = useNavigate();
    const [input, setInput] = useState("");
    const containerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (containerRef.current) {
            containerRef.current.scrollTop = containerRef.current.scrollHeight;
        }
    }, [messages]);

    const handleSend = (e: any) => {
        e.preventDefault();
        const text = input.trim();
        if (text) {
            onSendMessage(text);
            setInput("");
        }
    };

    const keyStr = localStorage.getItem("message-key");
    const ivStr = localStorage.getItem("message-iv");

    const key = keyStr ? CryptoJS.enc.Base64.parse(keyStr) : null;
    const iv = ivStr ? CryptoJS.enc.Base64.parse(ivStr) : null;

    const sortedMessages = useMemo(() => {
        return [...messages].sort(
            (a, b) => new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
        );
    }, [messages]);

    const decryptIfNeeded = (content: string): string => {
        if (!key || !iv) return content;
        try {
            if (/^[A-Za-z0-9+/=_-]+$/.test(content) && content.length > 16) {
                const plain = aesDecrypt(content, key, iv);
                return plain || content;
            }
            return content;
        } catch {
            return content;
        }
    };

    return (
        <Wrapper>
            <Header>
                <Avatar src={chatPhoto ?? "https://p7.hiclipart.com/preview/802/535/682/users-group-computer-icons-membership.jpg"} />
                <Username>{chatName}</Username>
                <span>{users.length + 1} учасників</span>
            </Header>

            <Container ref={containerRef}>
                {sortedMessages.length > 0 ? (
                    sortedMessages.map((msg) => {
                        const isCurrentUser = msg.senderId === currentUserId;
                        const decryptedText = decryptIfNeeded(msg.content);

                        return (
                            <MessageRow key={msg.id} isCurrentUser={isCurrentUser}>
                                {!isCurrentUser && (
                                    <AvatarChat
                                        onClick={() => navigate(ROUTES.USER_PROFILE.replace(":id", msg.sender.id))}
                                        src={msg.sender.profilePictureUrl}
                                        alt={msg.sender.username}
                                    />
                                )}
                                <MessageBubble isCurrentUser={isCurrentUser}>
                                    <UsernameChat
                                        onClick={() => navigate(ROUTES.USER_PROFILE.replace(":id", msg.sender.id))}
                                    >
                                        {msg.sender.username}
                                    </UsernameChat>
                                    <Content>{decryptedText}</Content>

                                    <Timestamp>{msg.createdAt.slice(11, 16)}</Timestamp>
                                </MessageBubble>

                                {isCurrentUser && (
                                    <AvatarChat
                                        onClick={() => navigate(ROUTES.USER_PROFILE.replace(":id", msg.sender.id))}
                                        src={msg.sender.profilePictureUrl}
                                        alt={msg.sender.username}
                                    />
                                )}
                            </MessageRow>
                        );
                    })
                ) : (
                    <p style={{ textAlign: "center", color: "#888", marginTop: "20px" }}>
                        Повідомлень ще немає
                    </p>
                )}
            </Container>

            <form onSubmit={handleSend}>
                <InputWrapper>
                    <Input
                        placeholder="Написати повідомлення..."
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                    />
                    <SendButton type="submit">
                        <img src={sendIcon} alt="send" />
                    </SendButton>
                </InputWrapper>
            </form>
        </Wrapper>
    );
};
