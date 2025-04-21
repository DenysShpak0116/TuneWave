import { FC, useState, KeyboardEvent, ChangeEvent } from "react";
import { Input, Tag, Wrapper } from "./multi-input.style";

interface MultiInputProps {
    name: string;
    placeholder: string;
    value: string[];
    onChange: (values: string[]) => void;
}

export const MultiInput: FC<MultiInputProps> = ({ placeholder, value, onChange }) => {
    const [inputValue, setInputValue] = useState("");

    const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter" && inputValue.trim() !== "") {
            e.preventDefault();
            if (!value.includes(inputValue.trim())) {
                onChange([...value, inputValue.trim()]);
            }
            setInputValue("");
        }
    };

    const removeTag = (tag: string) => {
        onChange(value.filter((t) => t !== tag));
    };

    return (
        <Wrapper>
            {value.map((tag) => (
                <Tag key={tag} onClick={() => removeTag(tag)}>
                    {tag} ✕
                </Tag>
            ))}
            <Input
                type="text"
                value={inputValue}
                placeholder={placeholder}
                onChange={(e: ChangeEvent<HTMLInputElement>) => setInputValue(e.target.value)}
                onKeyDown={handleKeyDown}
            />
        </Wrapper>
    );
};
