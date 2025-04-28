export interface IComment {
    id: string;
    author: {
        id: string;
        username: string
        profilePictureUrl: string;
    };
    content: string;
    createdAt: string;
}