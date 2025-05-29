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
  align-items: flex-start;
  padding-top: 50px;
  animation: ${fadeIn} 0.3s ease;
  z-index: 999;
`;

export const ModalContent = styled.div<{ $active: boolean }>`
  background: ${COLORS.dark_main};
  width: 100%;
  max-width: 600px;
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

export const Header = styled.div`
  display: flex;
  gap: 12px;
  flex-wrap: wrap;

  @media (max-width: 480px) {
    flex-direction: column;
    gap: 8px;
  }
`;

export const SearchInput = styled.input`
  flex: 1;
  padding: 10px 14px;
  font-size: 16px;
  border: none;
  border-radius: 8px;
  background-color: ${COLORS.dark_backdrop};
  color: ${COLORS.white};

  &::placeholder {
    color: ${COLORS.dark_secondary};
  }
`;

export const FilterButton = styled.button`
  padding: 10px 16px;
  font-size: 16px;
  background-color: ${COLORS.dark_backdrop};
  color: ${COLORS.white};
  border: none;
  border-radius: 8px;
  cursor: pointer;
  white-space: nowrap;

  &:hover {
    background-color: ${COLORS.dark_focusing};
  }
`;

export const ScrollContainer = styled.div`
  max-height: 400px;
  overflow-y: auto;
  padding-right: 10px;

  scrollbar-width: thin;
  scrollbar-color: #888 transparent;

  &::-webkit-scrollbar {
    width: 8px;
  }

  &::-webkit-scrollbar-thumb {
    background-color: #888;
    border-radius: 8px;
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }
`;

export const SongList = styled.ul`
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: .1rem;
`;

export const FilterContainer = styled.div`
  margin: 10px 0;
  display: grid;
  grid-template-columns: repeat(2,1fr);
  justify-items: center;
  padding: 10px;
  background: ${COLORS.dark_backdrop};
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
`;

export const FilterLabel = styled.label`
  display: block;
  font-size: 14px;
  color: ${COLORS.white};
  margin-bottom: 10px;
  font-weight: 500;
`;

export const SelectInput = styled.select`
  width: 200px;
  padding: 8px 12px;
  font-size: 14px;
  border-radius: 6px;
  border: 1px solid ${COLORS.dark_secondary};
  background-color: ${COLORS.dark_backdrop};
  color: ${COLORS.white};
  cursor: pointer;
  margin-top: 5px;

  &:focus {
    outline: none;
    border-color: ${COLORS.dark_additional};
  }
`;
