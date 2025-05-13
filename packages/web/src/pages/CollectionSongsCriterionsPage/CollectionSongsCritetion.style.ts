import styled from "styled-components";

export const Table = styled.table`
    width: 100%;
    border-collapse: collapse;
    margin-top: 24px;
`;

export const Th = styled.th`
    padding: 12px;
    background-color: #222;
    color: white;
    border: 1px solid #444;
    text-align: left;
`;

export const Td = styled.td`
    padding: 10px;
    border: 1px solid #ccc;
    text-align: left;
`;

export const DeleteButton = styled.button`
    padding: 6px 10px;
    background-color: #ff4d4f;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;

    &:hover {
        background-color: #d9363e;
    }
`;