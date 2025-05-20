import { MainLayout } from "@ui/layout/main-layout";
import { FC, FormEvent, useState } from "react";
import { Button, Container, Form, Input, Title } from "pages/ForgotPasswordPage/forgot-password.style";
import toast from "react-hot-toast";


export const ResetPasswordPage: FC = () => {
    const [newPassword, setNewPassword] = useState<string>("")
    const [repeatPassword, setRepeatPassword] = useState<string>("")
    const [token, setToken] = useState<string>("")

    const handleSubmit = (e: FormEvent) => {
        e.preventDefault()
        if (newPassword !== repeatPassword) {
            toast.error("Паролі не співпадають")
        }

    }

    return (
        <MainLayout>
            <Container>
                <Title>Відновлення паролю</Title>
                <Form onSubmit={handleSubmit}>
                    <Input
                        type="text"
                        placeholder="Введіть код відновлення"
                        value={token}
                        onChange={(e) => setToken(e.target.value)}
                        required
                    />
                    <Input
                        type="password"
                        placeholder="Введіть новий пароль"
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        required
                    />
                    <Input
                        type="password"
                        placeholder="Повторіть пароль"
                        value={repeatPassword}
                        onChange={(e) => setRepeatPassword(e.target.value)}
                        required
                    />
                    <Button type="submit">
                        Змінити пароль
                    </Button>
                </Form>
            </Container>
        </MainLayout>
    )
}