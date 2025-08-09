import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const UpdateUserPhotoContainer = styled.div`
    display: flex;
    justify-content: space-between;
    max-width: 400px;
    min-height: 150px;
    border-radius: 10px;
    padding: 20px;
    background-color: ${COLORS.dark_main};
`

export const UserImageContainer = styled.img`
    width: 150px;
    height: 150px;
    border-radius: 10px;
`

export const UploadBox = styled.div`
    height: 150px;
    border: 2px dashed ${COLORS.dark_additional};
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.2s ease;

    &:hover {
        border-color: ${COLORS.white};
    }
`;

export const PreviewImage = styled.img`
    max-width: 100%;
    max-height: 100%;
    border-radius: 6px;
    object-fit: cover;
`;

export const HiddenInput = styled.input`
    display: none;
`;

export const UploadIcon = styled.div`
    color: ${COLORS.dark_additional};
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    text-align: center;
    width: 200px;
    height: 200px;
    
    svg {
        width: 48px;
        height: 48px;
    }
     
    p {
        margin-top: 5px
    }
`;

export const UploadInformationContainer = styled.div`
    padding: 20px;
    max-width: 400px;
    flex-direction: column;
    margin-top: 24px;
    border-radius: 10px;
    background-color: ${COLORS.dark_main};

`
