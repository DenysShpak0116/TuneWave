type AuthInputType = { placeholder: string; type: string, name: string };

export const loginInputs: AuthInputType[] = [
    { placeholder: "Електрона пошта", type: "email", name: "email" },
    { placeholder: "Пароль", type: "password", name: "password" },
];
