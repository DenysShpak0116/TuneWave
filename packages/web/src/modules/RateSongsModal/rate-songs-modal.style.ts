import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Overlay = styled.div<{ $active: boolean }>`
    display: ${({ $active }) => ($active ? 'flex' : 'none')};
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    background: rgba(0, 0, 0, 0.5);
    justify-content: center;
    align-items: center;
    z-index: 1000;
`;

export const ModalContent = styled.div<{ $active: boolean; $animating: boolean }>`
    background: ${COLORS.dark_main};
    padding: 32px;
    border-radius: 16px;
    width: 600px;
    max-width: 90%;
    transition: transform 0.3s ease, opacity 0.3s ease;
    transform: ${({ $animating }) => ($animating ? 'translateY(20px)' : 'translateY(0)')};
    opacity: ${({ $animating }) => ($animating ? 0.5 : 1)};
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
`;


export const SongBlock = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 24px;
    margin-bottom: 32px;

    .song-card {
        flex: 1;
        display: flex;
        flex-direction: column;
        align-items: center;
        background: ${COLORS.dark_backdrop};
        padding: 16px;
        border-radius: 12px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        transition: transform 0.3s ease, opacity 0.3s ease;

        img {
            width: 120px;
            height: 120px;
            object-fit: cover;
            border-radius: 8px;
            margin-bottom: 4px;
        }

        .title {
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 4px;
            text-align: center;
        }

        .authors {
            font-size: 14px;
            color: ${COLORS.dark_additional};
            margin-bottom: 4px;
            text-align: center;
        }

        .genre,
        .duration {
            font-size: 13px;
            color: ${COLORS.dark_additional};
            margin-bottom: 4px;
        }

        .stats {
            display: flex;
            gap: 12px;
            font-size: 12px;
            margin-top: 4px;
            color: #555;

            span {
                display: flex;
                align-items: center;
                gap: 4px;
            }
        }

        p{
            margin-bottom: 4px;
            font-size: 14px;
            color: ${COLORS.dark_additional};
            text-align: center;
        }
    }
`;
export const ValuesBlock = styled.div`
    display: flex;
    flex-direction: row;
    width: 100%;
    justify-content: space-between;
`

export const SliderBlock = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 24px;
    width: 100%;

    input[type='range'] {
        width: 100%;
        margin: 16px 0;
    }

    span {
        display: inline-block;
        width: 20%;
        text-align: center;
        font-size: 14px;
        color: #666;
    }

    & > span:nth-child(1) { align-self: flex-start; }
    & > span:nth-child(2) { align-self: center; margin-left: -40%; }
    & > span:nth-child(3) { align-self: center; }
    & > span:nth-child(4) { align-self: center; margin-left: 40%; }
    & > span:nth-child(5) { align-self: flex-end; }
`;

export const NextButton = styled.button`
    padding: 12px 24px;
    background: ${COLORS.dark_main};
    color: ${COLORS.white};
    border: none;
    border-radius: 8px;
    font-size: 16px;
    cursor: pointer;
    transition: background 0.3s;

    &:hover {
        background: ${COLORS.dark_focusing}
    }

    &:disabled {
        background: #a5b4fc;
        cursor: not-allowed;
    }
`;


export const ComparisonTable = styled.table`
    width: 100%;
    border-collapse: collapse;
    margin-top: 16px;

    th, td {
        border: 1px solid #ccc;
        text-align: center;
        padding: 8px;
        font-size: 16px;
    }

    th {
        background-color: ${COLORS.dark_backdrop};
    }
`;

export const Leaderboard = styled.ul`
    list-style: none;
    padding: 0;
    margin: 24px 0;

    li {
        display: flex;
        align-items: center;
        padding: 8px 12px;
        margin-bottom: 8px;
        background: ${COLORS.dark_backdrop};
        border-radius: 8px;
        font-size: 16px;
        font-weight: 500;

    }
    .medal {
        margin-right: 12px;
        font-size: 20px;
    }
`;