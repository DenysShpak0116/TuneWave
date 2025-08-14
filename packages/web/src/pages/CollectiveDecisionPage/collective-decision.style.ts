import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const PageWrapper = styled.div`
`;

export const SectionTitle = styled.h2`
    font-size: 1.8rem;
    margin: 1rem 0 1rem;
    color: ${COLORS.white};
`;

export const SongList = styled.ul`
    list-style: none;
    padding: 0;
    margin: 0 0 2rem;
`;

export const SongItem = styled.li`
    padding: 0.5rem 0;
    font-size: 1.2rem;
    border-bottom: 1px solid #e0e0e0;
`;

export const ProfileTable = styled.table`
    width: 100%;
    border-collapse: collapse;
`;

export const TableRow = styled.tr`
    border-bottom: 1px solid #ddd;
`;

export const TableCell = styled.th`
    padding: 1rem;
    text-align: left;
    font-weight: 600;
    background-color: ${COLORS.dark_backdrop};
    border: 1px solid #ddd;
`;

export const SongListCell = styled.td`
    padding: 0.75rem;
    border: 1px solid #eee;
`;

export const UserName = styled.div`
    font-weight: 500;
    color: ${COLORS.dark_secondary};
`;