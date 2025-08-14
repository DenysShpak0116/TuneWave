import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const CardContainer = styled.div<{ isSelected: boolean }>`
  display: flex;
  width: 100%;
  align-items: flex-start;
  padding: 5px;
  border-radius: 5px;
  background-color: ${({ isSelected }) =>
        isSelected ? COLORS.dark_focusing : COLORS.chat_main};
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: ${({ isSelected }) =>
        isSelected ? COLORS.dark_focusing : COLORS.chat_focus};
  }
`;

export const Avatar = styled.img`
  width: 60px;
  height: 60px;
  border-radius: 8px;
  margin-right: 10px;
`;

export const Content = styled.div`
  display: flex;
  flex-direction: column;
`;

export const Username = styled.div`
  font-weight: bold;
  color: ${COLORS.white};
  margin-bottom: 4px;
`;

export const Message = styled.div`
  color: ${COLORS.dark_additional};
  font-size: 14px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
  max-height: calc(1.4em * 2);
`;