import { FC, ReactNode } from "react";
import { MainLayout } from "../main-layout";
import styled from "styled-components";
import { COLORS } from "@consts/colors.consts";

interface LayoutProps {
    children: ReactNode
}

const Container = styled.div`
  width: 100%;
  background-color: ${COLORS.dark_main};
  height: 575px;
  border-radius: 5px;
  margin: 0 auto;
  padding: 0 10px;
  display: grid;
  grid-template-columns: 1fr 2px 3fr;
  gap: 20px;
`

export const ChatLayout: FC<LayoutProps> = ({ children }) => {
    return (
        <MainLayout>
            <Container>{children}</Container>
        </MainLayout>
    )
}