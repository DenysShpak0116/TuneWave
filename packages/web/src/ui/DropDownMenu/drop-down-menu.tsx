import { FC, useState, useRef, useEffect } from "react"
import styled from "styled-components"
import moreIcon from "@assets/images/ic_more.png"
import { COLORS } from "@consts/colors.consts"

const MenuWrapper = styled.div`
  position: relative;
  display: flex;
`

const MenuButton = styled.button`
  background: transparent;
  border: none;
  cursor: pointer;
  justify-content: center;
`

const Dropdown = styled.div`
  position: absolute;
  top: 35px;
  background: ${COLORS.chat_main};
  border-radius: 8px;
  padding: 10px;
  min-width: 160px;
  z-index: 999;
`

interface DropdownMenuProps {
    children: React.ReactNode
}

export const DropdownMenu: FC<DropdownMenuProps> = ({ children }) => {
    const [open, setOpen] = useState(false)
    const ref = useRef<HTMLDivElement>(null)

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (ref.current && !ref.current.contains(event.target as Node)) {
                setOpen(false)
            }
        }

        document.addEventListener("mousedown", handleClickOutside)
        return () => document.removeEventListener("mousedown", handleClickOutside)
    }, [])

    return (
        <MenuWrapper ref={ref}>
            <MenuButton onClick={() => setOpen(!open)}>
                <img src={moreIcon} alt="more" width={10} />
            </MenuButton>
            {open && <Dropdown>{children}</Dropdown>}
        </MenuWrapper>
    )
}
