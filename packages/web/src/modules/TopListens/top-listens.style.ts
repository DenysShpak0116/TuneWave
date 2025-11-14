import styled, { keyframes, css } from "styled-components";
import { COLORS } from "@consts/colors.consts";

const fadeIn = keyframes`
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
`;

export const Container = styled.div`
  background: ${COLORS.dark_main};
  border-radius: 20px;
  padding: 28px;
  max-width: 1100px;
  margin: 0 auto;
  box-shadow: 0 10px 25px rgba(0,0,0,0.3);
  margin-bottom: 24px;
`;

export const Title = styled.h2`
  text-align: center;
  color: ${COLORS.white};
  margin-bottom: 32px;
  font-size: 1.8rem;
`;

export const LoaderWrapper = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
`;

export const Empty = styled.p`
  text-align: center;
  color: ${COLORS.white};
`;

export const Content = styled.div`
  display: grid;
  grid-template-columns: 1fr 340px;
  gap: 24px;

  @media (max-width: 900px) {
    grid-template-columns: 1fr;
  }
`;

export const TrackList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const TrackCard = styled.div<{ rank: number }>`
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 18px;
  border-radius: 14px;
  animation: ${fadeIn} 0.4s ease both;
  position: relative;
  background: ${COLORS.dark_backdrop};
  transition: transform 0.25s ease, box-shadow 0.25s ease;

  ${({ rank }) =>
    rank === 1 &&
    css`
      background: linear-gradient(135deg, #ffd700, #ffb347);
    `}
  ${({ rank }) =>
    rank === 2 &&
    css`
      background: linear-gradient(135deg, #dcdcdc, #a9a9a9);
    `}
  ${({ rank }) =>
    rank === 3 &&
    css`
      background: linear-gradient(135deg, #cd7f32, #8b4513);
    `}

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 18px rgba(0, 0, 0, 0.25);
  }
`;

export const RankBadge = styled.div<{ rank: number }>`
  font-size: 1.2rem;
  font-weight: bold;
  color: ${({ rank }) => (rank <= 3 ? "#222" : COLORS.white)};
  width: 40px;
  text-align: center;
`;

export const CoverWrapper = styled.div`
  position: relative;
  flex-shrink: 0;
`;

export const Overlay = styled.div`
  position: absolute;
  inset: 0;
  background: rgba(0,0,0,0.15);
  border-radius: 8px;
`;

export const Cover = styled.img`
  width: 68px;
  height: 68px;
  border-radius: 8px;
  object-fit: cover;
`;

export const Info = styled.div`
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

export const TrackTitle = styled.span`
  font-size: 1.1rem;
  font-weight: 600;
  color: ${COLORS.white};
`;

export const Genre = styled.span`
  font-size: 0.9rem;
  color: ${COLORS.white};
`;

export const Duration = styled.span`
  font-size: 0.85rem;
  color: ${COLORS.white};
  opacity: 0.8;
`;

export const Count = styled.span`
  font-weight: 600;
  font-size: 0.95rem;
  color: ${COLORS.dark_focusing};
  min-width: 130px;
  text-align: right;
`;

export const ChartCard = styled.div`
  background: ${COLORS.dark_backdrop};
  border-radius: 16px;
  padding: 16px;
  text-align: center;
  box-shadow: 0 4px 12px rgba(0,0,0,0.2);
  h3 {
    color: ${COLORS.white};
    margin-bottom: 16px;
  }
`;

export const LeftColumn = styled.div`
  display: flex;
  flex-direction: column;
  gap: 24px;
  flex: 1;
`;

export const ChartWrapper = styled.div`
  width: 100%;
  height: 250px;
`;

export const PopularContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 280px;

  @media (max-width: 900px) {
    min-width: 100%;
  }
`;

export const PopularList = styled.div`
  display: flex;
  flex-direction: column;
  gap: 12px;
`;

export const PopularCard = styled.div<{ isTop?: boolean }>`
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 10px;
  border-radius: 12px;
  background: ${({ isTop }) => (isTop ? "linear-gradient(135deg, #ffd700, #ffb347)" : COLORS.dark_backdrop)};
  cursor: pointer;
  position: relative;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  animation: ${fadeIn} 0.4s ease both;

  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 6px 14px rgba(0, 0, 0, 0.25);
  }
`;

export const TopBadge = styled.span`
  position: absolute;
  top: -6px;
  right: -6px;
  background: #ff4d4f;
  color: white;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 8px;
`;

export const PopularCoverWrapper = styled.div`
  position: relative;
  flex-shrink: 0;
`;

export const PopularCover = styled.img`
  width: 60px;
  height: 60px;
  border-radius: 10px;
  object-fit: cover;
`;

export const PopularInfo = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2px;
`;

export const PopularTrackTitle = styled.span`
  font-size: 1rem;
  font-weight: 600;
  color: ${COLORS.white};
`;

export const PopularGenre = styled.span`
  font-size: 0.85rem;
  color: ${COLORS.white};
`;

export const PopularDuration = styled.span`
  font-size: 0.8rem;
  color: ${COLORS.white};
  opacity: 0.8;
`;

export const PopularCount = styled.span`
  font-size: 0.85rem;
  color: ${COLORS.dark_focusing};
  font-weight: 600;
  margin-left: auto;
`;