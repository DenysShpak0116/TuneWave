import { COLORS } from "@consts/colors.consts";
import styled, { keyframes } from "styled-components";

const riseUp = keyframes`
  from {
    transform: translateY(30px) scale(0.95);
  }
  to {
    transform: translateY(0) scale(1);
  }
`;

export const Container = styled.div`
  width: 100%;
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
  background: ${COLORS.dark_main};
  border-radius: 16px;
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.05);
`;

export const Title = styled.h2`
  text-align: center;
  margin-bottom: 24px;
  color: ${COLORS.white};
`;

export const ChartWrapper = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 24px;
`;

export const ChartCard = styled.div`
  background: ${COLORS.dark_backdrop};
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);

  h3 {
    margin-bottom: 12px;
    font-size: 1.1rem;
    color: ${COLORS.white};
    text-align: center;
  }
`;

export const Message = styled.p`
  text-align: center;
  color: ${COLORS.white};
`;

export const Total = styled.p`
  text-align: center;
  margin-top: 12px;
  color: ${COLORS.white};
  font-weight: 600;
`;

export const StatsRow = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 32px;
`;

export const StatCard = styled.div`
  background: ${COLORS.dark_backdrop};
  border-radius: 12px;
  padding: 16px;
  text-align: center;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  animation: ${riseUp} 0.6s cubic-bezier(0.22, 1, 0.36, 1) forwards;
`;

export const StatValue = styled.div`
  font-size: 1.6rem;
  font-weight: 700;
  color: ${COLORS.white};
  margin-bottom: 4px;
`;

export const StatLabel = styled.div`
  font-size: 0.95rem;
  color: ${COLORS.white};
`;