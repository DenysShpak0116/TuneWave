import styled from "styled-components";
import { COLORS } from "@consts/colors.consts";

export const Wrapper = styled.div`
    padding: 0 24px;
`;

export const SongCardsContainer = styled.div`
    display: grid;
    grid-template-columns: repeat(5, 1fr);
    column-gap: 28px;
    margin-top: 16px;
`;

export const SongsText = styled.h3`
    margin: 0 auto;
    font-size: 32px;
    text-align: center;
    width: 100%;
`;

export const SongCard = styled.div`
    display: flex;
    flex-direction: column;
    background-color: ${COLORS.dark_main};
    padding: 16px;
    border-radius: 12px;
    transition: 
        transform 0.35s ease,
        background-color 0.35s ease,
        box-shadow 0.35s ease;

    &:hover {
        cursor: pointer;
        background-color: ${COLORS.dark_focusing};
        
    }

    &:hover .play-icon {
        opacity: 1;
        transform: translate(-50%, -50%) scale(1);
    }
`;

export const ImageWrapper = styled.div`
    position: relative;
    width: 100%;
    aspect-ratio: 1;
    border-radius: 8px;
    overflow: hidden;
`;

export const SongImage = styled.img`
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 8px;
    display: block;
`;

export const PlayIcon = styled.img`
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) scale(0.8);
    width: 40px;
    height: 40px;
    opacity: 0;
    transition: all 0.3s ease;
    pointer-events: auto;
    z-index: 2;
`;

export const SongTitle = styled.div`
    margin-top: 12px;
    font-weight: 600;
    font-size: 16px;
    color: white;
`;

export const SongArtist = styled.div`
    padding-top: 5px;
    font-size: 10px;
    color: ${COLORS.dark_additional};
`;
