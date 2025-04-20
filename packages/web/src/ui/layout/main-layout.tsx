import { FC, ReactNode } from "react"
import { Container, Wrapper } from "./main-layout.style"
import { Header } from "@components/Header/header.component"


interface MainLayoutProps {
    children: ReactNode
}

export const MainLayout: FC<MainLayoutProps> = ({children}) => {

    return(
        <>
        <Wrapper>
            <Header/>
            <Container>{children}</Container>
        </Wrapper>
        </>
    )
}