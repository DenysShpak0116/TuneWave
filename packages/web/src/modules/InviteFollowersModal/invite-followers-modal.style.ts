import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const ListContainer = styled.div`
  margin-top: 15px;
  overflow-y: auto;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const FollowerItem = styled.div<{ $selected: boolean }>`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  background: ${({ $selected }) => ($selected ? COLORS.dark_focusing : COLORS.dark_backdrop)};
  transition: 0.2s;

  &:hover {
    background: ${({ $selected }) => ($selected ? COLORS.dark_focusing : COLORS.dark_backdrop)};
  }
`;

export const Avatar = styled.img`
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
`;

export const Username = styled.span`
  font-size: 16px;
  font-weight: 500;
`;

export const ChatNameInput = styled.input`
  width: 50%;
  padding: 10px;
  margin-top: 15px;
  border-radius: 8px;
  border: 1px solid ${COLORS.dark_backdrop};
  background-color: ${COLORS.dark_backdrop};
  color: ${COLORS.dark_additional};
  font-size: 16px;
  outline: none;
  transition: 0.2s;

  &:focus {
    border-color: ${COLORS.dark_focusing};
  }
`;

export const CreateButton = styled.button`
  margin-top: 15px;
  padding: 10px;
  border: none;
  border-radius: 8px;
  background: ${COLORS.dark_focusing};
  color: white;
  font-size: 16px;
  cursor: pointer;
  transition: 0.2s;

  &:disabled {
    background: ${COLORS.dark_backdrop};
    cursor: not-allowed;
  }

  &:not(:disabled):hover {
    background: ${COLORS.dark_backdrop};
  }
`