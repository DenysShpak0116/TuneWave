import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: ${COLORS.dark_main};
`;

export const Header = styled.div`
    padding: 12px 20px;
    margin-top: 10px;
    background-color: ${COLORS.chat_main};
    display: flex;
    align-items: center;
    gap: 12px;
    border-radius: 5px;
`;

export const Avatar = styled.img`
    width: 32px;
    height: 32px;
    border-radius: 50%;
`;

export const Username = styled.div`
    color: white;
    font-weight: 500;
`;

export const Container = styled.div`
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 20px;
  background-color: transparent;
  max-height: 400px;

  scrollbar-width: thin;
  scrollbar-color: #444 transparent;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: #444;
    border-radius: 4px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }
`;

export const MessageRow = styled.div<{ isCurrentUser: boolean }>`
    display: flex;
    justify-content: ${({ isCurrentUser }) =>
        isCurrentUser ? "flex-end" : "flex-start"};
`;

export const MessageBubble = styled.div<{ isCurrentUser: boolean }>`
  background-color: ${({ isCurrentUser }) =>
        isCurrentUser ? COLORS.dark_focusing : COLORS.chat_user_message};
  color: white;
  border-radius: ${({ isCurrentUser }) =>
        isCurrentUser ? "20px 20px 1px" : "20px 20px 20px 1px"};
  padding: 10px 14px;
  max-width: 300px;
  word-wrap: break-word;
  font-size: 14px;
  display: inline-flex;
  align-items: flex-end;
  gap: 6px;
`;

export const Timestamp = styled.div`
  font-size: 10px;
  color: rgba(255, 255, 255, 0.5);
  white-space: nowrap;
  flex-shrink: 0;
`;

export const InputWrapper = styled.div`
  padding: 15px 20px;
  display: flex;
  align-items: center;
  border-radius: 12px;
  background-color: transparent;
  flex-shrink: 0;
`;
export const Input = styled.input`
    flex-grow: 1;
    padding: 10px 15px;
    border-radius: 12px;
    border: none;
    background-color: ${COLORS.chat_main};
    color: white;
    font-size: 14px;

    &::placeholder {
        color: #999;
    }

    &:focus {
        outline: none;
    }
`;

export const SendButton = styled.button`
    background: transparent;
    border: none;
    margin-left: 10px;
    cursor: pointer;

    img {
        width: 20px;
        height: 20px;
    }
`;