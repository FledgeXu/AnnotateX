import { useNavigate } from "@tanstack/react-router";
import { Eye, EyeOff } from "lucide-react";
import { toast } from "sonner";
import { useForm } from "react-hook-form";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { createAPI } from "@/config";
import { store } from "@/store";
import { useStoreActions } from "easy-peasy";
import type { StoreModel } from "@/store/types";
import { useState } from "react";
import { useMutation } from "@tanstack/react-query";
import type { Response, LoginToken } from "@/models";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { cn } from "@/lib/utils";

const schema = z.object({
  username: z.string().min(3, "Username is required"),
  password: z.string().min(6, "Password must be at least 6 characters"),
  agree: z
    .boolean()
    .refine((val) => val === true, { message: "You must agree to the terms" }),
});

type FormData = z.infer<typeof schema>;

export const login = async (data: FormData): Promise<Response<LoginToken>> => {
  const api = createAPI(store);
  const res = await api.post("/v1/auth/login", {
    username: data.username,
    password: data.password,
  });
  return res.data;
};

export const useLoginMutation = () => {
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
      const message = error?.response?.data?.message || "Failed to log in.";
      toast.error(`Login failed: ${message}`);
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

export const LoginForm = () => {
  const form = useLoginForm();
  const [visible, setVisible] = useState(false);
  const loginMutation = useLoginMutation();
  const onSubmit = form.handleSubmit((data) => loginMutation.mutate(data));
  return (
    <>
      <Form {...form}>
        <form onSubmit={onSubmit}>
          <div className="flex flex-col gap-4">
            <FormField
              control={form.control}
              name="username"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Username</FormLabel>
                  <FormControl>
                    <Input placeholder="Enter your account name." {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="password"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Password</FormLabel>
                  <FormControl>
                    <div className="relative">
                      <Input
                        type={cn(visible ? "text" : "password")}
                        id="password"
                        placeholder="Enter your password"
                        className="pr-10"
                        {...field}
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7"
                        onClick={() => setVisible((prev) => !prev)}
                      >
                        {visible ? (
                          <EyeOff className="w-4 h-4" />
                        ) : (
                          <Eye className="w-4 h-4" />
                        )}
                        <span className="sr-only">
                          Toggle password visibility
                        </span>
                      </Button>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="agree"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <div className="flex items-center gap-3">
                      <Checkbox
                        id="terms"
                        checked={field.value}
                        onCheckedChange={(val) => field.onChange(val === true)}
                      />
                      <Label htmlFor="terms">
                        I have read and agree to{" "}
                        <a href="/terms" className="underline">
                          Terms of Service{" "}
                        </a>
                      </Label>
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="w-full">
              Login
            </Button>
          </div>
        </form>
      </Form>
    </>
  );
};
