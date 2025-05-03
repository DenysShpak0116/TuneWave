import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Container = styled.div`
    margin-top: 10px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    background-color: ${COLORS.dark_backdrop};
    border-radius: 10px;
    cursor: pointer;
`;

export const Image = styled.img`
    width: 50px;
    height: 50px;
    object-fit: cover;
    border-radius: 8px;
    margin-right: 16px;
`;

export const Info = styled.div`
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
`;

export const Title = styled.span`
    font-size: 16px;
    font-weight: 600;
    color: white;
`;

export const Description = styled.span`
    font-size: 14px;
    color: #aaa;
`;

export const StyledRadio = styled.input`
  appearance: none;
  width: 18px;
  height: 18px;
  border: 2px solid ${COLORS.dark_main};
  border-radius: 50%;
  position: relative;
  transition: 0.2s;
  pointer-events: none;

  &:checked::before {
    content: '';
    position: absolute;
    top: 4px;
    left: 4px;
    width: 8px;
    height: 8px;
    background-color: ${COLORS.dark_focusing};
    border-radius: 50%;
  }
  &:focus{
    background-color: ${COLORS.dark_focusing};
  }
`;