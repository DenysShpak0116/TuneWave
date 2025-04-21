export type SongDetailsInputType = {
    placeholder: string
    type: 'text'
    name: string;
    isMulti: boolean
}

export const songDetailsInputs: SongDetailsInputType[] = [
    { placeholder: "Назва", type: "text", name: "title", isMulti: false },
    { placeholder: "Жанри", type: "text", name: "genre", isMulti: false },
    { placeholder: "Артисти", type: "text", name: "artists", isMulti: true },
    { placeholder: "Теги", type: "text", name: "tags", isMulti: true },
]