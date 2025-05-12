import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const CollectionSongsContainer = styled.div`
  width: 100%;
  border-radius: 10px;
  background-color: ${COLORS.dark_main};
  padding: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
`;

export const TableHeader = styled.div`
  display: grid;
  grid-template-columns: 40px 1fr 1fr 140px 80px 40px;
  padding: 12px 16px;
  font-size: 14px;
  color: ${COLORS.dark_secondary};
  background-color: ${COLORS.dark_main};
  border-bottom: 1px solid ${COLORS.dark_additional};
`;

export const TableRow = styled.div<{ active?: boolean }>`
  display: grid;
  grid-template-columns: 40px 1fr 1fr 140px 80px 40px;
  align-items: center;
  padding: 10px 16px;
  background-color: ${({ active }) =>
    active ? COLORS.dark_focusing : COLORS.dark_backdrop};
  border-radius: 8px;
  margin: 6px 8px;
  transition: background-color 0.2s ease-in-out;

`;

export const IndexBox = styled.div`
  text-align: center;
  color: ${COLORS.dark_secondary};
`;

export const CoverAndInfo = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

export const Cover = styled.img`
  width: 40px;
  height: 40px;
  border-radius: 6px;
`;

export const SongTextInfo = styled.div`
  display: flex;
  flex-direction: column;
`;

export const Title = styled.div`
  font-size: 15px;
  color: ${COLORS.dark_secondary};
  font-weight: 500;
`;

export const Author = styled.div`
  font-size: 12px;
  color: ${COLORS.dark_secondary};
`;

export const Album = styled.div`
  font-size: 14px;
  color: ${COLORS.dark_secondary};
`;

export const DateAdded = styled.div`
  font-size: 14px;
  color: ${COLORS.dark_secondary};
`;

export const Duration = styled.div`
  font-size: 14px;
  color: ${COLORS.dark_secondary};
`;

export const Options = styled.div`
  text-align: center;
  color: ${COLORS.dark_secondary};
  font-size: 18px;
  cursor: pointer;
`;