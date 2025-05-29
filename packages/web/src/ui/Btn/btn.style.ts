import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const StyledButton = styled.button`
  background-color: ${COLORS.dark_main};
  color: #fff;
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s ease;

  &:hover {
    background-color: ${COLORS.dark_focusing}
  }

  &:active {
    transform: scale(0.98);
  }
`;