import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { createAPI } from "@/config";
import { store } from "@/store";
import { useStoreActions } from "easy-peasy";
import type { StoreModel } from "@/store/types";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import type { Response, LoginToken } from "@/models";
import {
    CardHeader,
    CardTitle,
    Card,
    CardContent,
    CardFooter,
    CardDescription,
} from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";

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

export const login = async (data: FormData): Promise<Response<LoginToken>> => {
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
            authLogin(data.data.token);
            navigate({ to: "/" });
        },
        onError: (error: any) => {
            const msg = error?.response?.data?.message || "Failed to log in.";
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
    const [visible, setVisible] = useState(false);
    const loginMutation = useLoginMutation(setLoginErrorMessage);
    const onSubmit = handleSubmit((data) => loginMutation.mutate(data));

    return (
        <div className="flex items-center justify-center h-screen">
            <Card className="w-full max-w-sm">
                <CardHeader>
                    <CardTitle>
                        <Label className="text-xl font-bold">Login</Label>
                    </CardTitle>
                </CardHeader>
                <CardContent>
                    <form>
                        <div className="flex flex-col gap-6">
                            <div className="grid gap-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    type="email"
                                    placeholder="m@example.com"
                                    required
                                />
                            </div>
                            <div className="grid gap-2">
                                <div className="flex items-center">
                                    <Label htmlFor="password">Password</Label>
                                    <a
                                        href="#"
                                        className="ml-auto inline-block text-sm underline-offset-4 hover:underline"
                                    >
                                        Forgot your password?
                                    </a>
                                </div>
                                <Input id="password" type="password" required />
                            </div>
                            <div className="grid gap-2">
                                <div className=" flex items-center space-x-2">
                                    <Checkbox id="terms" />
                                    <Label htmlFor="terms">Accept terms and conditions</Label>
                                </div>
                            </div>
                        </div>
                    </form>
                </CardContent>
                <CardFooter className="flex-col gap-2">
                    <Button type="submit" className="w-full">
                        Login
                    </Button>
                    <CardDescription>
                        Don't have account?
                        <a href="#" className="underline ml-2">
                            Register!
                        </a>
                    </CardDescription>
                </CardFooter>
            </Card>
        </div>
    );
}
