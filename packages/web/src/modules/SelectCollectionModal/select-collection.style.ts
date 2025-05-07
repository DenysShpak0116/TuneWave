import { COLORS } from "@consts/colors.consts";
import styled, { keyframes } from "styled-components";

const fadeIn = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`;

const slideUp = keyframes`
  from { transform: translateY(50px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
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

export const ModalHeader = styled.div`
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    width: 100%;
    text-align: center;
    color: ${COLORS.white};
`;

export const ModalHeaderText = styled.h3`
    justify-content: center;
    font-size: 20px;
    font-weight: 600;
    margin-bottom: 5px;
`

export const AddCollectionIcon = styled.img`
    width: 24px;
    height: 24px;
    cursor: pointer;
`

export const ModalContent = styled.div<{ $active: boolean }>`
    background: ${COLORS.dark_main};
    width: 600px;
    height: 500px;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    animation: ${slideUp} 0.3s ease;
    padding: 20px;
    box-sizing: border-box;
`;

export const ModalBody = styled.div`
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex: 1;
  max-height: 375px;
  overflow-y: auto;
  padding-right: 8px;

      &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: #ccc;
    border-radius: 10px;
    transition: background 0.3s ease;
  }

  &::-webkit-scrollbar-thumb:hover {
    background-color: #999;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }
`;