import { COLORS } from "@consts/colors.consts";
import styled, { keyframes } from "styled-components";

const fadeIn = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`;

const slideDown = keyframes`
  from {
    transform: translateY(-50px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
`;

export const Overlay = styled.div<{ $active: boolean }>`
    display: ${({ $active }) => ($active ? "flex" : "none")};
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    justify-content: center;
    align-items: center;
    animation: ${fadeIn} 0.3s ease;
    z-index: 999;
`;

export const ModalContent = styled.div<{ $active: boolean }>`
  background: ${COLORS.dark_main};
  width: 100%;
  max-width: 300px;
  height: auto;
  max-height: 90vh;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  animation: ${slideDown} 0.3s ease;
  padding: 20px;
  box-sizing: border-box;
  overflow-y: auto;

  @media (max-width: 768px) {
    padding: 16px;
    border-radius: 10px;
  }

  @media (max-width: 480px) {
    padding: 12px;
    border-radius: 8px;
  }
`;