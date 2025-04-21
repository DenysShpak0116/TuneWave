import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Wrapper = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 0 24px;
`;

export const UserBlock = styled.div`
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    border-radius: 10px;
    width: 100%;
    background-color: ${COLORS.dark_main};
    padding: 24px;
    position: relative;
`;

export const Avatar = styled.img`
    width: 100px;
    height: 100px;
    border-radius: 10px;
    object-fit: cover;
`;

export const Name = styled.h2`
    font-size: 24px;
    margin: 0;
`;

export const Stats = styled.div`
    margin-top: 4px;
    font-size: 14px;
    color: ${COLORS.dark_additional};
`;

export const Bio = styled.p`
    margin-top: 8px;
    font-size: 14px;
    color: ${COLORS.white};
`;

export const SettingsIcon = styled.div`
    position: absolute;
    top: 20px;
    right: 20px;
    cursor: pointer;
    color: ${COLORS.white};
`;