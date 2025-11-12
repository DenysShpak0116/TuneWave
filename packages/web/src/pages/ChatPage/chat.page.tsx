// chat.page.tsx
import { FC, useEffect, useMemo, useRef, useState } from "react";
import { ChatLayout } from "@ui/layout/Chat/chat-layout";
import { Loader } from "@ui/Loader/loader.component";
import { ChatList } from "@modules/ChatList";
import { Divider } from "./chat.style";
import { MainChat } from "@modules/MainChat";
import { useChatSocket } from "@modules/MainChat/hooks/useChatSocket"; // new version (string key)
import { useGetUserChats } from "./hooks/useGetUserChats";
import { IMessageType } from "types/chat/message.type";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { useLocation, useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

export const ChatPage: FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const searchParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
    const currentUserId = useAuthStore(state => state.user!.id)
    const rawUserIds = searchParams.get("userIds") ?? "";

    const userIds = useMemo(() => {
        if (!rawUserIds) return [];
        return rawUserIds
            .split(",")
            .filter(Boolean)
            .filter(id => id !== currentUserId);
    }, [rawUserIds, currentUserId]);

    const chatName = searchParams.get("name") ?? "";
    const { data: chatPreviews, isLoading, refetch } = useGetUserChats();
    const [messages, setMessages] = useState<IMessageType[]>([]);
    const sortedUserIds = useMemo(() => [...userIds].sort(), [userIds]);
    const usersParam = useMemo(() => sortedUserIds.join(","), [sortedUserIds]);
    const chatKey = useMemo(() => `${usersParam}_${chatName}`, [usersParam, chatName]);

    const currentChatPreview = useMemo(() => {
        if (!chatPreviews) return undefined;
        return chatPreviews.find(cp => {
            const cpUsersSorted = [...cp.userIds].sort().join(",");
            if (chatName) {
                return cpUsersSorted === usersParam && cp.chatName === chatName;
            }
            return cpUsersSorted === usersParam;
        });
    }, [chatPreviews, usersParam, chatName]);

    const { sendMessage } = useChatSocket({
        usersParam,
        chatName,
        onMessageReceived: (msg: IMessageType) => {
            setMessages(prev => [...prev, msg]);
        }
    });

    const prevChatKey = useRef<string>("");
    useEffect(() => {
        if (prevChatKey.current !== chatKey) {
            setMessages([]);
            prevChatKey.current = chatKey;
        }
    }, [chatKey]);

    const didRedirect = useRef(false);
    useEffect(() => {
        if (userIds.length === 0 && chatPreviews && chatPreviews.length > 0 && !didRedirect.current) {
            didRedirect.current = true;
            const firstChat = chatPreviews[0];
            navigate(
                `${ROUTES.CHAT_PAGE}?userIds=${firstChat.userIds.join(",")}&name=${encodeURIComponent(firstChat.chatName ?? "")}`,
                { replace: true }
            );
        }
    }, [userIds, chatPreviews, navigate]);

    const handleSendMessage = (text: string) => {
        sendMessage(text);
        refetch();
    };

    if (isLoading) {
        return (
            <ChatLayout>
                <Loader />
            </ChatLayout>
        );
    }

    const noChats = (!chatPreviews || chatPreviews.length === 0) && userIds.length === 0;
    return (
        <ChatLayout>
            <ChatList chatPreviews={chatPreviews} targetUserIds={userIds} />
            <Divider />

            {noChats ? (
                <div style={{ padding: "20px", fontSize: "16px" }}>
                    Ви поки що не маєте чатів
                </div>
            ) : usersParam === "" ? (
                <div style={{ padding: "20px", fontSize: "16px" }}>
                    Виберіть чат
                </div>
            ) : (
                <MainChat
                    chatName={chatName || currentChatPreview?.chatName || "Чат"}
                    messages={messages}
                    currentUserId={currentUserId}
                    users={sortedUserIds}
                    chatPhoto={"https://cdn-icons-png.flaticon.com/512/681/681494.png"}
                    onSendMessage={handleSendMessage}
                    chatId={currentChatPreview?.id}
                />
            )}
        </ChatLayout>
    );
};
