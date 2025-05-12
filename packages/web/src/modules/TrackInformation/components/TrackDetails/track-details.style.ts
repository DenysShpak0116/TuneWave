import styled from "styled-components";
import { COLORS } from "@consts/colors.consts";
import { Link } from "react-router-dom";

export const TrackDetailsContainer = styled.div`
  display: flex;
  flex-direction: column;
  padding: 20px;
  width: 210px;
  background-color: ${COLORS.dark_main};
  border-radius: 10px;
  color: ${COLORS.white};
  font-family: "Inter", sans-serif;
  font-size: 14px;
  gap: 5px;
  grid-column: 1;
  grid-row: 2;
  align-self: start;
`;

export const TrackInfoBlock = styled.div`
    display: flex;
    flex-direction: column;
    margin-top: 5px;
`;

export const TrackInfoTitle = styled.span`
    font-weight: 600;
    color: ${COLORS.dark_additional};
`;

export const TrackInfoText = styled.span`
    color: ${COLORS.dark_secondary};
    font-weight: 400;
`;

export const UserInfoText = styled(Link)`
    color: ${COLORS.dark_secondary};
    font-weight: 400;
    text-decoration: underline;
`
