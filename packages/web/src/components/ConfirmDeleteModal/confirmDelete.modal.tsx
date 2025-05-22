import { FC } from "react";
import { ModalContent, Overlay } from "./confirmDelete.style";
import { Button } from "@ui/Btn/btn.component";

interface IConfirmDeleteModal {
    active: boolean;
    setActive: (value: boolean) => void;
    onDelete: () => void
}

export const ConfirmDeleteModal: FC<IConfirmDeleteModal> = ({ active, setActive, onDelete }) => {

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} onClick={(e) => e.stopPropagation()}>
                <p style={{ textAlign: "center" }}>Ви впевнені, що хочете видалити пісню</p>
                <div style={{ display: "flex", justifyContent: "space-around", marginTop: "14px" }}>
                    <Button text="Так" onClick={onDelete} />
                    <Button text="Ні" onClick={() => setActive(false)} />
                </div>
            </ModalContent>
        </Overlay>
    )
}