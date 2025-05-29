type AuthInputType = { placeholder: string; type: string, name: string };

export const registrationInputs: AuthInputType[] = [
    { placeholder: "Ім’я користувача", type: "text", name: "username" },
    { placeholder: "Електрона пошта", type: "email", name: "email" },
    { placeholder: "Пароль", type: "password", name: "password" },
    { placeholder: "Повторити пароль", type: "password", name: "repeatPassword" },
];

