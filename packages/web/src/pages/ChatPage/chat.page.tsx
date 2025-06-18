import { FC, useEffect, useState } from "react";
import { ChatLayout } from "@ui/layout/Chat/chat-layout";
import { Loader } from "@ui/Loader/loader.component";
import { ChatList } from "@modules/ChatList";
import { Divider } from "./chat.style";
import { MainChat } from "@modules/MainChat";
import { useChatSocket } from "@modules/MainChat/hooks/useChatSocket";
import { useGetUserChats } from "./hooks/useGetUserChats";
import { IMessageType } from "types/chat/message.type";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { useGetUser } from "pages/UserProfilePage/hooks/useGetUserById";
import { useLocation, useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

export const ChatPage: FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const searchParams = new URLSearchParams(location.search);
    const targetUserId = searchParams.get("targetUserId");

    const { data: chatPreviews, isLoading, refetch } = useGetUserChats();
    const [messages, setMessages] = useState<IMessageType[]>([]);
    const currentUserId = useAuthStore(state => state.user!.id);

    const shouldFetchUser = !!targetUserId;
    const { data: user, isLoading: loadUser } = useGetUser(targetUserId!)

    const { sendMessage } = useChatSocket({
        userId: targetUserId ?? "",
        onMessageReceived: (msg) => {
            setMessages((prev) => [...prev, msg]);
        }
    });

    useEffect(() => {
        setMessages([]);
    }, [targetUserId]);

    useEffect(() => {
        if (!targetUserId && chatPreviews && chatPreviews.length > 0) {
            const firstChat = chatPreviews[0];
            navigate(`${ROUTES.CHAT_PAGE}?targetUserId=${firstChat.targetUserId}`, { replace: true });
        }
    }, [targetUserId, chatPreviews, navigate]);

    const handleSendMessage = (text: string) => {
        sendMessage(text);
        refetch()
    };

    if (isLoading || (shouldFetchUser && loadUser)) {
        return (
            <ChatLayout>
                <Loader />
            </ChatLayout>
        );
    }

    return (
        <ChatLayout>
            <ChatList chatPreviews={chatPreviews} targetUserId={targetUserId ?? ""} />
            <Divider />
            {chatPreviews?.length === 0 && !targetUserId || !user ? (
                <div style={{ padding: "20px", fontSize: "16px" }}>
                    Ви поки що не маєте чатів
                </div>
            ) : (
                <MainChat
                    targetUserId={user.id}
                    messages={messages}
                    currentUserId={currentUserId}
                    partnerUsername={user.username}
                    partnerAvatar={user.profilePictureUrl}
                    onSendMessage={handleSendMessage}
                />
            )}
        </ChatLayout>
    );
};
