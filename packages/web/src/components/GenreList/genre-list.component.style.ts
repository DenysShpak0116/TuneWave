import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";


interface CardProps {
  background: string;
}

export const Grid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
`;

export const Card = styled.div<CardProps>`
  position: relative;
  height: 180px;
  background: ${(props) => props.background};
  border-radius: 12px;
  padding: 16px;
  overflow: hidden;
  transition: background 0.3s ease;
  cursor: pointer;
`;

export const GenreTitle = styled.h3`
  color: #ffffff;
  font-size: 24px;
  font-weight: bold;
  z-index: 2;
  position: relative;
`;

export const CoverWrapper = styled.div`
  position: absolute;
  bottom: 0;
  right: 0;
  width: 80px;
  transform: rotate(15deg) translate(-55px, 10px);
`;

export const CoverImage = styled.img`
  width: 150px;
  height: 150px;
  border-radius: 8px;
  display: block;
`;
