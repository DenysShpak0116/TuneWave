import { FC, useRef, useState, ChangeEvent } from "react";
import {
    ModalContent,
    Overlay,
    RightSide,
    ModalHeader,
    ModalBody,
    ModalHeaderText,
    BackIcon
} from "./collection-modal.style";

import { InputField } from "@ui/Input/input-field.component";
import { TextAreaField } from "@ui/TextArea/text-area-field.component";
import {
    PreviewImage,
    UploadBox,
    UploadContainer,
    UploadIcon
} from "@components/ContentContainer/content-container.style";
import { FiUpload } from "react-icons/fi";
import { Button } from "@ui/Btn/btn.component";
import { useCreateCollection } from "./hooks/useCreateCollection";
import backIcon from "@assets/images/ic_arrow.png"
import toast from "react-hot-toast";

interface ICollectionModal {
    active: boolean;
    setActive: (value: boolean) => void;
    userId: string
    onCreated?: () => void;
}

export const CollectionModal: FC<ICollectionModal> = ({ active, setActive, userId, onCreated }) => {
    const [collectionName, setCollectionName] = useState<string>("");
    const [collectionDescription, setCollectionDescription] = useState<string>("");
    const [coverFile, setCoverFile] = useState<File | null>(null);
    const inputRef = useRef<HTMLInputElement | null>(null);
    const [coverPreview, setCoverPreview] = useState<string | null>(null);
    const { mutate: submitCollection } = useCreateCollection();

    const handleCoverUpload = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setCoverFile(file);
            const reader = new FileReader();
            reader.onloadend = () => {
                setCoverPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSubmit = () => {
        if (!coverFile) {
            return;
        }

        const formData = new FormData();
        formData.append("userId", userId);
        formData.append("title", collectionName);
        formData.append("description", collectionDescription);
        formData.append("cover", coverFile);

        submitCollection(formData, {
            onSuccess: () => {
                toast.success("Колекцію успішно створено!");
                onCreated?.();
                setActive(false);
                resetForm();
            }
        });

    };

    const resetForm = () => {
        setActive(false)
        setCollectionName("")
        setCollectionDescription("")
        setCoverFile(null)
        setCoverPreview(null)
    }

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} onClick={(e) => e.stopPropagation()}>
                <ModalHeader>
                    <BackIcon src={backIcon} onClick={() => setActive(false)} />
                    <ModalHeaderText>
                        Створити колекцію
                    </ModalHeaderText>
                </ModalHeader>
                <ModalBody>
                    <RightSide>
                        <InputField
                            placeholder="Назва колекції"
                            label=""
                            value={collectionName}
                            onChange={setCollectionName}
                        />

                        <TextAreaField
                            placeholder="Опис колекції"
                            label=""
                            value={collectionDescription}
                            onChange={setCollectionDescription}
                        />
                    </RightSide>

                    <UploadContainer>
                        <UploadBox onClick={() => inputRef.current?.click()}>
                            {coverPreview ? (
                                <PreviewImage src={coverPreview} alt="preview" />
                            ) : (
                                <UploadIcon>
                                    <FiUpload size={36} />
                                </UploadIcon>
                            )}
                            <input
                                type="file"
                                accept="image/*"
                                onChange={handleCoverUpload}
                                ref={inputRef}
                                style={{ display: "none" }}
                            />
                        </UploadBox>
                    </UploadContainer>
                </ModalBody>
                <Button text={"Додати колекцію"} onClick={handleSubmit} />
            </ModalContent>
        </Overlay>
    );
};
