import { FC } from "react";
import { IComment } from "types/comments/comment.type";
import { CommentAuthor, CommentContent, CommentHeader, CommentItem, CommentText, CommentUserAvatar, DeleteText } from "./comment.style";
import { parseDate } from "helpers/date-parse";
import { ROUTES } from "pages/router/consts/routes.const";
import { useNavigate } from "react-router-dom";


interface ICommentProps {
    comment: IComment;
    userId: string | undefined
    deleteCommentFn: (id: string) => void
}
export const Comment: FC<ICommentProps> = ({ comment, userId, deleteCommentFn }) => {
    const navigate = useNavigate()

    return (
        <CommentItem>
            <CommentUserAvatar src={comment.author.profilePictureUrl} onClick={() => navigate(ROUTES.USER_PROFILE.replace(":id", comment.author.id))} />
            <CommentContent>
                <CommentHeader>
                    <CommentAuthor>{comment.author.username}</CommentAuthor>
                    <span>{parseDate(comment.createdAt)}</span>
                </CommentHeader>
                <CommentText>{comment.content}</CommentText>
                {comment.author.id === userId && (
                    <DeleteText onClick={() => deleteCommentFn(comment.id)}>Delete</DeleteText>
                )}
            </CommentContent>
        </CommentItem>
    )
}