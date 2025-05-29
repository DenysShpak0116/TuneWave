import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const CommentItem = styled.div`
  display: flex;
  gap: 12px;
  background-color: transparent;
  border-radius: 8px;
  padding: 10px 14px;
`;

export const CommentUserAvatar = styled.img`
  width: 40px;
  height: 40px;
  border-radius: 25%;
  object-fit: cover;
  cursor: pointer;
`;

export const CommentContent = styled.div`
  display: flex;
  flex-direction: column;
  flex: 1;
`;

export const CommentHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: ${COLORS.dark_additional};
`;

export const CommentAuthor = styled.span`
  font-weight: bold;
`;

export const DeleteText = styled.span`
  margin-top: 5px;
  font-size: 11px;
  cursor: pointer;

  &:hover{
    text-decoration: underline;
  }
`


export const CommentText = styled.p`
  color: ${COLORS.white};
  font-size: 14px;
  margin-top: 4px;
`;