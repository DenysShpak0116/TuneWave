import styled from "styled-components";

export const Wrapper = styled.div`
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 8px;
    background-color: transparent;
`;

export const Tag = styled.span`
    background-color: #333;
    color: #fff;
    padding: 4px 8px;
    border-radius: 6px;
    font-size: 14px;
    cursor: pointer;
`;

export const Input = styled.input`
    border: none;
    background: transparent;
    color: white;
    font-size: 14px;
    outline: none;
    min-width: 100px;
    flex-grow: 1;
`;
