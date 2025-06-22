import { createFileRoute } from "@tanstack/react-router";
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
import { useRouter } from "@tanstack/react-router";

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

function Login() {
    const {
        control,
        handleSubmit,
        formState: { errors },
    } = useForm<FormData>({
        resolver: zodResolver(schema),
        defaultValues: {
            username: "",
            password: "",
            agree: false,
        },
    });

    const authLogin = useStoreActions<StoreModel>(
        (actions) => actions.auth.login,
    );
    const router = useRouter();

    const onSubmit = async (data: FormData) => {
        const api = createAPI(store);
        try {
            const res = await api.post("/v1/auth/login", {
                username: data.username,
                password: data.password,
            });
            const { token } = res.data;
            authLogin(token);
            router.navigate({ to: "/" });
        } catch (err: any) {
            console.error("Fail", err.message);
        }
    };

    return (
        <Flex align="center" justify="center" height="100vh">
            <Card size="3" className="w-full max-w-sm">
                <form onSubmit={handleSubmit(onSubmit)}>
                    <Flex gap={"3"} direction={"column"}>
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
                        {Object.entries(errors).map(([key, error]) => (
                            <Callout.Root key={key} color="red">
                                <Callout.Icon>
                                    <InfoCircledIcon />
                                </Callout.Icon>
                                <Callout.Text>{error?.message}</Callout.Text>
                            </Callout.Root>
                        ))}
                    </Flex>
                </form>
            </Card>
        </Flex>
    );
}
