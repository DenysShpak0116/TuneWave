import { COLORS } from "@consts/colors.consts"
import styled from "styled-components"

export const Container = styled.div`
  max-width: 600px;
  margin: 0 auto;
  padding: 2rem;
`

export const Title = styled.h2`
  text-align: center;
  margin-bottom: 1.5rem;
`
export const Input = styled.input`
  width: 100%;
  background-color: transparent;
  padding: 0.5rem;
  margin-bottom: 1rem;
  border-radius: 6px;
  color: ${COLORS.dark_secondary};
  border: 1px solid ${COLORS.dark_secondary};
`

export const Button = styled.button`
  padding: 0.5rem 1rem;
  background-color: ${COLORS.dark_main};
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 1.5rem;

  &:hover {
    background-color: ${COLORS.dark_focusing};
  }

  &:disabled {
    background-color: ${COLORS.dark_main};
    cursor: not-allowed;
  }
`



export const List = styled.ul`
  list-style: none;
  padding: 0;
`

export const ListItem = styled.li`
  display: flex;
  flex-direction: row;
  background: ${COLORS.dark_main};
  margin-bottom: 0.5rem;
  padding: 0.75rem 1rem;
  border-radius: 6px;
`

export const DeleteButton = styled.button`
  background-color:rgb(225, 84, 81);
  color: white;
  border: none;
  border-radius: 6px;
  padding: 0.25rem 0.5rem;
  margin-left: auto;
  cursor: pointer;
  justify-self: end;

  &:hover {
    background-color: #d32f2f;
  }
`
export const EditButton = styled(DeleteButton)`
    background: #ffc107;

    &:hover {
    background-color: #ffc107;
  }
`

export const SaveButton = styled(DeleteButton)`
    background: #28a745;
    
    &:hover {
    background-color: #28a745;
  }
`
