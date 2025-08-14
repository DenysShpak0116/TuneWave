import { COLORS } from "@consts/colors.consts";
import styled, { css, keyframes } from "styled-components";


const fadeIn = keyframes`
  from { opacity: 0; }
  to { opacity: 1; }
`;
const fadeOut = keyframes`
  from { opacity: 1; transform: translateY(0); }
  to { opacity: 0; transform: translateY(30px); }
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
  align-items: flex-start;
  padding-top: 50px;
  animation: ${fadeIn} 0.3s ease;
  z-index: 999;
`;

export const ModalContent = styled.div<{ $active: boolean; $animating: boolean }>`
  background: ${COLORS.dark_main};
  max-width: 600px;
  width: 100%;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  padding: 20px;
  animation: ${({ $animating }) =>
    $animating ? css`${fadeOut} 0.3s ease forwards` : css`${slideUp} 0.3s ease`};
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

export const Header = styled.div`
  display: flex;
  gap: 12px;
  flex-wrap: wrap;

  @media (max-width: 480px) {
    flex-direction: column;
    gap: 8px;
  }
`;

export const SongHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;

  p {
    justify-content: flex-end;
  }

  img {
    width: 200px;
    height: 200px;
    object-fit: cover;
    border-radius: 8px;
  }

  h3 {
    font-size: 24px;
    color: ${COLORS.white};
    margin: 0;
  }
`;

export const InputsBlock = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 10px;

  label {
    font-size: 14px;
    color: ${COLORS.dark_additional};
    margin-right: 5px;
  }

  input {
    padding: 6px 10px;
    border-radius: 6px;
    border: 1px solid ${COLORS.dark_secondary};
    background: transparent;
    color: ${COLORS.white};
  }
`;

export const NextButton = styled.button`
  margin-top: 20px;
  align-self: flex-end;
  background-color: ${COLORS.dark_focusing};
  color: ${COLORS.white};
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;

  &:hover {
    background-color: ${COLORS.dark_focusing};
  }
`;