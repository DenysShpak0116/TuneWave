import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const FormContainer = styled.div`
    padding: 20px;
    background-color: ${COLORS.dark_main};
    border-radius: 10px;
`;

export const Title = styled.h2`
    color: white;
    text-align: center;
    margin-bottom: 20px;
`;

export const Container = styled.div`
    display: flex;
    justify-content: space-around;
    flex-wrap: wrap;
`;

export const UploadContainer = styled.div`
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 10px;
`;

export const UploadBox = styled.div`
    width: 150px;
    height: 150px;
    border: 2px dashed gray;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 10px;
    cursor: pointer;
    overflow: hidden;
`;

export const UploadIcon = styled.div`
    color: gray;
`;

export const PreviewImage = styled.img`
    width: 100%;
    height: 100%;
    object-fit: cover;
`;

export const HiddenInput = styled.input`
    display: none;
`;

export const Divider = styled.hr`
    margin: 20px 0;
    border: 0.5px solid gray;
    width: 100%;
`;

export const InputsContainer = styled.div`
    display: flex;
    flex-direction: column;
    gap: 10px
`
