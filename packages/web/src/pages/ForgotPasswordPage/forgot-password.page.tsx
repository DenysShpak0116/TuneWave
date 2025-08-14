import { MainLayout } from "@ui/layout/main-layout";
import { FC, useState } from "react";
import {
    Container,
    Title,
    Form,
    Input,
    Button,
} from "./forgot-password.style";
import { useSendTokenToEmail } from "./hooks/useForgotPassword";

export const ForgotPasswordPage: FC = () => {
    const [email, setEmail] = useState("");
    const { mutate: sendToken, isLoading,} = useSendTokenToEmail();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        sendToken(email);
    };

    return (
        <MainLayout>
            <Container>
                <Title>Відновлення паролю</Title>
                <Form onSubmit={handleSubmit}>
                    <Input
                        type="email"
                        placeholder="Введіть свій email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                    <Button type="submit" disabled={isLoading}>
                        {isLoading ? "Відправка..." : "Відправити код відновлення"}
                    </Button>
                </Form>
            </Container>
        </MainLayout>
    );
};
