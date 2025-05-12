import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Card = styled.div`
  display: flex;
  background: ${COLORS.dark_backdrop};
  border-radius: 12px;
  margin-top: 10px;

  overflow: hidden;
  transition: background-color 0.4s ease;
  cursor: pointer;

  &:hover {
    background-color: ${COLORS.dark_focusing};
  }
`;

export const Cover = styled.img`
  width: 100px;
  height: 100px;
  object-fit: cover;
`;

export const Info = styled.div`
  flex: 1;
  padding: 0.75rem;
`;

export const Title = styled.h3`
  margin: 0;
  font-size: 1.1rem;
`;

export const Subtitle = styled.p`
  margin: 0.25rem 0;
  color: #666;
  font-size: 0.9rem;
`;

export const Artist = styled.p`
  margin: 0.25rem 0;
  font-size: 0.85rem;
  color: #333;
`;

export const Icon = styled.img`
    width: 14px;
    height: 14px;
`

export const Stats = styled.div`
  margin-top: 0.5rem;
  font-size: 0.85rem;
  display: flex;
  gap: 10px;
  color: #555;
`;
