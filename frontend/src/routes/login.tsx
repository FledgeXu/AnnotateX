import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
    Box,
    Button,
    Card,
    Checkbox,
    Flex,
    Heading,
    Link,
    Text,
    TextField,
    IconButton,
    Callout,
} from "@radix-ui/themes";
import { EyeOpenIcon, InfoCircledIcon } from "@radix-ui/react-icons";
import { z } from "zod";
import { createAPI } from "@/config";
import { store } from "@/store";
import { useStoreActions } from "easy-peasy";
import type { StoreModel } from "@/store/types";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";

const schema = z.object({
    username: z.string().min(3, "Username is required"),
    password: z.string().min(6, "Password must be at least 6 characters"),
    agree: z
        .boolean()
        .refine((val) => val === true, { message: "You must agree to the terms" }),
});

type FormData = z.infer<typeof schema>;

export const Route = createFileRoute("/login")({
    component: Login,
});

export const login = async (data: FormData): Promise<{ token: string }> => {
    const api = createAPI(store);
    const res = await api.post("/v1/auth/login", {
        username: data.username,
        password: data.password,
    });
    return res.data;
};

export const useLoginMutation = (onErrorMessage: (msg: string) => void) => {
    const authLogin = useStoreActions<StoreModel>(
        (actions) => actions.auth.login,
    );
    const navigate = useNavigate();

    return useMutation({
        mutationFn: login,
        onSuccess: (data) => {
            authLogin(data.token);
            navigate({ to: "/" });
        },
        onError: (error: any) => {
            const msg = error?.response?.data?.error || "Failed to log in.";
            onErrorMessage(msg);
        },
    });
};

const useLoginForm = () => {
    return useForm<FormData>({
        resolver: zodResolver(schema),
        defaultValues: {
            username: "",
            password: "",
            agree: false,
        },
    });
};

function Login() {
    const {
        control,
        handleSubmit,
        formState: { errors },
    } = useLoginForm();
    const [loginErrorMessage, setLoginErrorMessage] = useState<string | null>(
        null,
    );
    const loginMutation = useLoginMutation(setLoginErrorMessage);
    const onSubmit = handleSubmit((data) => loginMutation.mutate(data));

    return (
        <Flex align="center" justify="center" height="100vh">
            <Card size="3" className="w-full max-w-sm">
                <form onSubmit={onSubmit}>
                    <Flex gap={"3"} direction={"column"}>
                        {Object.entries(errors).map(([key, error]) => (
                            <Callout.Root key={key} color="red">
                                <Callout.Icon>
                                    <InfoCircledIcon />
                                </Callout.Icon>
                                <Callout.Text>{error?.message}</Callout.Text>
                            </Callout.Root>
                        ))}
                        {loginErrorMessage && (
                            <Callout.Root key="loginErrorMessage" color="red">
                                <Callout.Icon>
                                    <InfoCircledIcon />
                                </Callout.Icon>
                                <Callout.Text>{loginErrorMessage}</Callout.Text>
                            </Callout.Root>
                        )}
                        <Heading>Login</Heading>
                        <Controller
                            name="username"
                            control={control}
                            render={({ field }) => (
                                <Box>
                                    <Text weight={"medium"}>Username</Text>
                                    <TextField.Root
                                        type="text"
                                        placeholder="Enter your account name."
                                        {...field}
                                    >
                                        <TextField.Slot />
                                    </TextField.Root>
                                </Box>
                            )}
                        />
                        <Box>
                            <Text weight={"medium"}>Password</Text>

                            <Controller
                                name="password"
                                control={control}
                                render={({ field }) => (
                                    <TextField.Root
                                        type="password"
                                        placeholder="Enter your password."
                                        {...field}
                                    >
                                        <TextField.Slot />
                                        <TextField.Slot>
                                            <IconButton variant="ghost">
                                                <EyeOpenIcon />
                                            </IconButton>
                                        </TextField.Slot>
                                    </TextField.Root>
                                )}
                            />
                        </Box>
                        <Flex align="center" gap="2">
                            <Controller
                                name="agree"
                                control={control}
                                render={({ field }) => (
                                    <Checkbox
                                        checked={field.value}
                                        onCheckedChange={(val) => field.onChange(val === true)}
                                    />
                                )}
                            />
                            <Text size="2">
                                I have read and agree to{" "}
                                <Link href="/terms" target="_blank">
                                    Terms of Service
                                </Link>
                            </Text>
                        </Flex>
                        <Box>
                            <Flex direction={"row-reverse"}>
                                <Button type="submit">Login</Button>
                            </Flex>
                        </Box>
                    </Flex>
                </form>
            </Card>
        </Flex>
    );
}
