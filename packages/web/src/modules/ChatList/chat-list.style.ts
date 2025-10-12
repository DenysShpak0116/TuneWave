import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";


export const ChatPreviewContainer = styled.div`
    display: flex;
    flex-direction: column;
    gap: 8px;
`

export const SearchInput = styled.input`
    width: 100%;
    background-color: ${COLORS.dark_backdrop};
    border-radius: 5px;
    padding: 10px;
    color: ${COLORS.dark_additional}
`