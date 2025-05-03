import styled from "styled-components";

export const TrackInfoWrapper = styled.div`
  display: flex;
  align-items: center;
  gap: 12px;
`;

export const TrackLogo = styled.img`
  width: 40px;
  height: 40px;
  border-radius: 8px;
  object-fit: cover;
  cursor: pointer;
`;

export const TextBlock = styled.div`
  display: flex;
  flex-direction: column;
  overflow: hidden;
`;

export const TrackName = styled.span`
  font-size: 16px;
  font-weight: 600;
  color: white;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

export const TrackArtist = styled.span`
  font-size: 14px;
  color: #aaa;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
`;

export const AddIcon = styled.img`
  width: 24px;
  height: 24px;
  align-self: center;
`