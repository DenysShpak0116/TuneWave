import { useEffect, useRef } from "react";
import { IMessageType } from "types/chat/message.type";

interface UseChatSocketProps {
    userId: string;
    onMessageReceived: (message: IMessageType) => void;
}

export const useChatSocket = ({ userId, onMessageReceived }: UseChatSocketProps) => {
    const socketRef = useRef<WebSocket | null>(null);
    const authToken = localStorage.getItem("token");

    useEffect(() => {
        const url = `ws://localhost:8081/ws/chat?targetUserId=${userId}&authToken=${authToken}`;
        const socket = new WebSocket(url);
        socketRef.current = socket;

        socket.onopen = () => {
            console.log("🔗 WebSocket connected");
        };

        socket.onmessage = (event) => {
            const message: IMessageType = JSON.parse(event.data);
            onMessageReceived(message);
        };


        socket.onclose = (event) => {
            console.log("❌ WebSocket disconnected", event);
        };

        socket.onerror = (error) => {
            console.error("⚠️ WebSocket error:", error);
        };

        return () => {
            socket.close();
        };
    }, [authToken, userId]);

    const sendMessage = (content: string) => {
        if (socketRef.current?.readyState === WebSocket.OPEN) {
            const message = { content };
            socketRef.current.send(JSON.stringify(message));
        } else {
            console.warn("⛔ WebSocket is not open. Message not sent.");
        }
    };

    return { sendMessage };
};
