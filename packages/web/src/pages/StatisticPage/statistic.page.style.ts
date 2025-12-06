import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Nav = styled.nav`
  display: flex;
  justify-content: center;
  gap: 16px;
  margin-bottom: 24px;
`;

export const NavButton = styled.button<{ active: boolean }>`
  background: ${({ active }) => (active ? COLORS.dark_focusing : COLORS.dark_backdrop)};
  color: ${COLORS.white};
  border: none;
  border-radius: 8px;
  padding: 10px 18px;
  cursor: pointer;
  font-weight: 500;
  transition: 0.2s;

  &:hover {
    background: ${COLORS.dark_focusing};
  }
`;