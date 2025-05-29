import styled from "styled-components"
import { COLORS } from "@consts/colors.consts"

export const Wrapper = styled.div`
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
`

export const ButtonGroup = styled.div`
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
`

export const TabButton = styled.button<{ active: boolean }>`
  padding: 10px 16px;
  margin: 0 5px;
  border: none;
  border-radius: 20px;
  background-color: ${({ active }) => (active ? "#6441a5" : "#2e2e2e")};
  color: white;
  cursor: pointer;
`

export const UserListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 12px;
`

export const UserCard = styled.div`
  background: ${COLORS.dark_main};
  border-radius: 12px;
  padding: 10px 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
`

export const Avatar = styled.img`
  width: 40px;
  height: 40px;
  border-radius: 50%;
`

export const Username = styled.div`
  flex-grow: 1;
  font-size: 16px;
  color: white;
`

export const DeleteBtn = styled.button`
  background:rgb(236, 94, 120);
  padding: 6px 12px;
  color: white;
  border-radius: 8px;
  font-size: 14px;
  text-decoration: none;
`

export const EmptyText = styled.div`
  text-align: center;
  color: #999;
  padding: 20px;
`
