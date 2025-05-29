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
    flex-direction: column;
    text-align: center;
    width: 100%;
    text-align: center;
    color: ${COLORS.white};

`;

export const BackIcon = styled.img`
    width: 24px;
    height: 24px;
    cursor: pointer;
`

export const ModalHeaderText = styled.h3`
    font-size: 20px;
    font-weight: 600;
`

export const ModalContent = styled.div<{ $active: boolean }>`
    background: ${COLORS.dark_main};
    width: 600px;
    height: 400px;
    border-radius: 12px;
    display: flex;
    flex-direction: column;
    animation: ${slideUp} 0.3s ease;
    padding: 25px;
    box-sizing: border-box;
`;

export const ModalBody = styled.div`
    display: flex;
    flex-direction: row-reverse;
    justify-content: space-evenly;
    gap: 15px;
    flex: 1;
    margin-bottom: 10px;
`;


export const RightSide = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 10px;
`;

export const BottomSide = styled.div`
    display: flex;
    justify-content: center;
    padding-top: 20px;
`;

export const Input = styled.input`
    padding: 10px;
    font-size: 16px;
    border-radius: 8px;
    border: 1px solid #ccc;
`;

export const Textarea = styled.textarea`
    padding: 10px;
    font-size: 16px;
    border-radius: 8px;
    border: 1px solid #ccc;
    resize: none;
    height: 80px;
`;

export const UploadLabel = styled.div`
    padding: 10px 20px;
    background: #007bff;
    color: #fff;
    border-radius: 8px;
    cursor: pointer;
    font-weight: bold;
`;
