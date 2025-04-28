import { ChangeEvent, FC, useRef, useState } from "react";
import { IUser } from "types/user/user.type";
import { HiddenInput, UpdateUserPhotoContainer, UploadBox, UploadIcon, UploadInformationContainer, UserImageContainer } from "./update-user.style";
import { FiUpload } from "react-icons/fi";
import { InputField } from "@ui/Input/input-field.component";
import { TextAreaField } from "@ui/TextArea/text-area-field.component";
import { Button } from "@ui/Btn/btn.component";
import { useUpdateUser } from "./hooks/useUpdateUser";
import toast from "react-hot-toast";

interface IUpdateUserForm {
    user: IUser;
}

export const UpdateUserForm: FC<IUpdateUserForm> = ({ user }) => {
    const [coverPreview, setCoverPreview] = useState<string | null>(null);
    const [username, setUsername] = useState<string>(user.username);
    const [userInfo, setUserInfo] = useState<string>(user.profileInfo);

    const inputRef = useRef<HTMLInputElement>(null);

    const { mutate: updateUserInfo, isPending } = useUpdateUser();

    const handleCoverUpload = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onloadend = () => {
                setCoverPreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleSubmit = () => {
        updateUserInfo(
            { id: user.id, profileInfo: userInfo, username },
            {
                onSuccess: () => {
                    toast.success("Профіль оновлено!");
                },
                onError: () => {
                    toast.error("Помилка оновлення профілю");
                }
            }
        );
    };

    return (
        <>
            <UpdateUserPhotoContainer>
                <UserImageContainer src={coverPreview || user.profilePictureUrl} />
                <UploadBox onClick={() => inputRef.current?.click()}>
                    <UploadIcon>
                        <FiUpload size={16} />
                        <p>Завантажити нову обкладинку</p>
                    </UploadIcon>

                    <HiddenInput
                        type="file"
                        accept="image/*"
                        onChange={handleCoverUpload}
                        ref={inputRef}
                    />
                </UploadBox>
            </UpdateUserPhotoContainer>

            <UploadInformationContainer>
                <p>Інформація:</p>
                <InputField
                    label="Ім'я користувача"
                    value={username}
                    onChange={setUsername}
                    placeholder="Введіть ім'я користувача"
                />
                <TextAreaField
                    label="Опис користувача"
                    placeholder="Введіть опис користувача"
                    value={userInfo}
                    onChange={setUserInfo}
                />
                <Button
                    style={{ marginTop: "20px" }}
                    text={isPending ? "Збереження..." : "Зберегти зміни"}
                    onClick={handleSubmit}
                    isDisabled={isPending}
                />
            </UploadInformationContainer>
        </>
    );
};
