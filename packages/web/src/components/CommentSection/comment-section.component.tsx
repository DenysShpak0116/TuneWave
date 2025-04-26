import { FC, useState, KeyboardEvent } from "react";
import {
    CommentSectionContainer,
    CommentInput,
    CommentList,
    Avatar,
    StyledInputWrapper,
    SendIcon,
    NoContentText
} from "./comment-section.style";
import { IComment } from "types/comments/comment.type";
import { Comment } from "@ui/Comment/comment.component";
import sendIcon from "@assets/images/ic_send.png";
import { useCreateComment } from "./hooks/useCreateComment";
import toast from "react-hot-toast";
import { useDeleteComment } from "./hooks/useDeleteComment";

interface ICommentSectionProps {
    userAvatar: string | undefined;
    comments: IComment[];
    userId: string | undefined
    songId: string;
    onNewComment: (comment: IComment) => void;
    onDeleteComment: (id: string) => void;
}


export const CommentSection: FC<ICommentSectionProps> = ({ userAvatar, comments, userId, songId, onNewComment, onDeleteComment }) => {
    const [inputValue, setInputValue] = useState("");
    const { mutate: createComment } = useCreateComment();
    const { mutate: removeComment } = useDeleteComment();

    const handleSend = () => {
        if (inputValue.trim() && userId) {
            const content = inputValue;
            setInputValue("");
            createComment(
                { userId, songId, content },
                {
                    onSuccess: (newComment) => {
                        onNewComment(newComment);
                    },
                    onError: () => {
                        toast.error("Помилка відправки")
                    }
                }
            );
        }
    };

    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") {
            handleSend();
        }
    };

    const deleteComment = (id: string) => {
        removeComment(id, {
            onSuccess: () => {
                toast.success("Коментар видалено!");
                onDeleteComment(id)
            },
            onError: () => {
                toast.error("Помилка видалення коментаря");
            }
        });
    };

    console.log(userId)

    if (comments.length == 0) {
        return (
            <CommentSectionContainer>
                <NoContentText>Немає коментарів :( </NoContentText>
            </CommentSectionContainer>
        )
    }

    return (
        <CommentSectionContainer>
            {userId !== undefined && (
                <StyledInputWrapper>
                    <Avatar src={userAvatar} />
                    <CommentInput
                        placeholder="Напишіть ваш коментар"
                        value={inputValue}
                        onChange={(e) => setInputValue(e.target.value)}
                        onKeyDown={handleKeyDown}
                    />
                    <SendIcon src={sendIcon} onClick={handleSend} />
                </StyledInputWrapper>
            )}
            <CommentList>
                {comments.map(comment => (
                    <Comment
                        deleteCommentFn={deleteComment}
                        userId={userId}
                        key={comment.id}
                        comment={comment} />
                ))}
            </CommentList>
        </CommentSectionContainer>
    );
};
