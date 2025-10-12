// useChatSocket.ts
import { useEffect, useRef, useMemo } from "react";
import { IMessageType } from "types/chat/message.type";

interface UseChatSocketProps {
    usersParam: string;
    chatName: string;
    onMessageReceived: (message: IMessageType | { __reset: true }) => void;
}

export const useChatSocket = ({ usersParam, chatName, onMessageReceived }: UseChatSocketProps) => {
    const socketRef = useRef<WebSocket | null>(null);
    const prevChatKey = useRef<string>("");

    const authToken = localStorage.getItem("token");
    const chatKey = useMemo(() => `${usersParam}_${chatName}`, [usersParam, chatName]);

    const shouldConnect =
        !!authToken &&
        usersParam.trim() !== "" &&
        chatName.trim() !== "";

    useEffect(() => {
        if (prevChatKey.current !== "" && prevChatKey.current !== chatKey) {
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            onMessageReceived({ __reset: true } as any);
        }
        prevChatKey.current = chatKey;
    }, [chatKey]);

    useEffect(() => {
        if (!shouldConnect) return;

        const url = `ws://localhost:8081/ws/chat?authToken=${encodeURIComponent(
            authToken!
        )}&userIds=${encodeURIComponent(usersParam)}&name=${encodeURIComponent(chatName)}`;

        const socket = new WebSocket(url);
        socketRef.current = socket;

        socket.onopen = () => {
            console.log("🔗 WebSocket connected:", url);
        };

        socket.onmessage = (event) => {
            try {
                const message: IMessageType = JSON.parse(event.data);
                onMessageReceived(message);
            } catch (err) {
                console.error("Failed to parse message", err);
            }
        };

        socket.onclose = (event) => {
            console.log("❌ WebSocket disconnected", event);
        };

        socket.onerror = (error) => {
            console.error("⚠️ WebSocket error:", error);
        };

        return () => {
            try {
                socket.close();
            } catch (e) {
                console.error(e);
            }
            socketRef.current = null;
        };
    }, [shouldConnect, authToken, usersParam]); 

    const sendMessage = (content: string) => {
        if (socketRef.current?.readyState === WebSocket.OPEN) {
            socketRef.current.send(JSON.stringify({ content }));
        } else {
            console.warn("⛔ WebSocket is not open. Message not sent.");
        }
    };

    return { sendMessage };
};
