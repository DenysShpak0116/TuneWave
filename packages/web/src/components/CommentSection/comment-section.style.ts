import styled from "styled-components";
import { COLORS } from "@consts/colors.consts";

export const CommentSectionContainer = styled.div`
  width: 100%;
  grid-column: 2;
  grid-row: 2;
  border-radius: 10px;
  background-color: ${COLORS.dark_main};
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 300px;
`;

export const NoContentText = styled.p`
  font-size: 32px;
  justify-content: center;
  text-align: center;
`

export const StyledInputWrapper = styled.div`
  display: flex;
  align-items: center;
  background-color: ${COLORS.dark_main};
  border-radius: 10px;
  padding: 8px 12px;
  border: 1px solid ${COLORS.dark_additional};
`;

export const Avatar = styled.img`
  width: 32px;
  height: 32px;
  border-radius: 25%;
  margin-right: 12px;
`;

export const SendIcon = styled.img`
  width: 24px;
  height: 24px;
  justify-content: end;
  cursor: pointer;
`

export const CommentInput = styled.input`
  flex: 1;
  background: transparent;
  border: none;
  color: white;
  font-size: 14px;

  &:focus {
    outline: none;
  }

  &::placeholder {
    color: ${COLORS.dark_additional};
  }
`;

export const CommentList = styled.div`
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;

  scrollbar-width: thin;
  scrollbar-color: ${COLORS.dark_additional} transparent;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: ${COLORS.dark_additional};
    border-radius: 4px;
  }
`;