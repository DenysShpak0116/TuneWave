import styled, { keyframes } from "styled-components";
import { COLORS } from "@consts/colors.consts";

const slideUp = keyframes`
  from { transform: translateY(10px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
`;

export const Container = styled.div`
  width: 100%;
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px;
  background: ${COLORS.dark_main};
  border-radius: 16px;
`;

export const Title = styled.h2`
  text-align: center;
  margin-bottom: 24px;
  color: ${COLORS.white};
`;

export const Tabs = styled.div`
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-bottom: 24px;
`;

export const TabButton = styled.button<{ active: boolean }>`
  background: ${({ active }) => (active ? COLORS.dark_focusing : COLORS.dark_backdrop)};
  color: ${COLORS.white};
  border: none;
  border-radius: 10px;
  padding: 10px 20px;
  cursor: pointer;
  transition: background 0.25s, transform 0.15s;

  &:hover {
    transform: translateY(-2px);
    background: ${COLORS.dark_focusing};
  }
`;

export const ChartCard = styled.div`
  background: ${COLORS.dark_backdrop};
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  margin-bottom: 28px;
  animation: ${slideUp} 0.4s ease;
  h3 {
    color: ${COLORS.white};
    text-align: center;
    margin-bottom: 12px;
  }
`;

export const CombinationList = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
`;

export const ComboCard = styled.div`
  background: ${COLORS.dark_backdrop};
  border-radius: 12px;
  padding: 14px;
  animation: ${slideUp} 0.4s ease;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
`;

export const CoversRow = styled.div`
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  justify-content: center;
`;

export const TrackCover = styled.img`
  width: 64px;
  height: 64px;
  border-radius: 8px;
  object-fit: cover;
`;

export const TrackCoverSkeleton = styled.div`
  width: 64px;
  height: 64px;
  border-radius: 8px;
  background: ${COLORS.white};
`;

export const ComboFooter = styled.div`
  display: flex;
  justify-content: space-between;
  margin-top: 12px;
`;

export const ComboText = styled.span`
  color: ${COLORS.white};
  font-size: 0.9rem;
`;

export const ComboCount = styled.span`
  color: ${COLORS.dark_focusing};
  font-weight: 600;
`;

export const Empty = styled.p`
  text-align: center;
  color: ${COLORS.white};
`;

export const LoaderWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
`;

export const TrackPreviewWrapper = styled.div`
  position: relative;
  display: inline-block;
`;

export const TrackCount = styled.span`
  position: absolute;
  bottom: -4px;
  right: -4px;
  background: ${COLORS.dark_focusing};
  color: ${COLORS.white};
  font-size: 0.75rem;
  border-radius: 6px;
  padding: 2px 5px;
  line-height: 1;
`;
