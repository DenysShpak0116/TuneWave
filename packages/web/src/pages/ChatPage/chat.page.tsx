import { ChatList } from "@modules/ChatList";
import { ChatLayout } from "@ui/layout/Chat/chat-layout";
import { FC } from "react";
import { IChatPreviewType } from "types/chat/chat-preview";
import { Divider } from "./chat.style";
import { IMessageType } from "types/chat/message.type";
import { MainChat } from "@modules/MainChat";

const chatPreview: IChatPreviewType[] = [
    {
        id: "1233",
        userAvatar: 'https://i.pinimg.com/736x/13/8c/93/138c93cd2cf946e4a58c04d77c347fb6.jpg',
        username: "someMan",
        lastMessage: "Hey guy Hey guy Hey guy Hey guy "
    },
    {
        id: "1233121",
        userAvatar: 'https://i.pinimg.com/736x/13/8c/93/138c93cd2cf946e4a58c04d77c347fb6.jpg',
        username: "someMan",
        lastMessage: "wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup wazup "
    },
]

const messages: IMessageType[] = [
    {
        id: "1",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Worem ipsum dolor sit amet",
        chatId: "chat1",
        senderId: "user1",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
    {
        id: "2",
        createdAt: "2025-05-18T17:50:00Z",
        content: "Lorem ipsum dolor sit",
        chatId: "chat1",
        senderId: "me",
    },
];


export const ChatPage: FC = () => {
    return (
        <ChatLayout>
            <ChatList chatPreviews={chatPreview} />
            <Divider />
            <MainChat
                messages={messages}
                currentUserId="me"
                partnerUsername="guntersteam"
                partnerAvatar="https://i.pinimg.com/736x/13/8c/93/138c93cd2cf946e4a58c04d77c347fb6.jpg"
                onSendMessage={() =>console.log("sss")
                }

            />
        </ChatLayout>
    )
}