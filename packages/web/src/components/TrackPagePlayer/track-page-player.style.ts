import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const TrackPagePlayerContainer = styled.div`
  width: 100%;
  grid-column: 2;
  grid-row: 1;
  border-radius: 10px;
  background-color: ${COLORS.dark_main};
  padding: 16px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 16px;
  max-height: 300px;
`;

export const TrackTitle = styled.h1`
  font-size: 48;
  font-weight: 600;
  text-align: center;
`
