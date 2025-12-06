export interface IMessageType {
    id: string;
    createdAt: string;
    content: string;
    chatId: string;
    senderId: string;
    sender: ISender;
}

interface ISender {
    id: string;
    username: string;
    role: string
    profilePictureUrl: string
    profileInfo: string
    followers: number
}