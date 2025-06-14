import { createFileRoute } from "@tanstack/react-router";
import React, { useState } from "react";
import * as Checkbox from "@radix-ui/react-checkbox";
import * as Label from "@radix-ui/react-label";
import { CheckIcon, EyeClosedIcon, EyeOpenIcon } from "@radix-ui/react-icons";

export const Route = createFileRoute("/login")({
    component: Login,
});

function Login() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [agreed, setAgreed] = useState(false);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!agreed) {
            alert("请先同意用户协议。");
            return;
        }
        alert(`欢迎你，${username}！`);
    };

    return (
        <form onSubmit={handleSubmit} >
            <h2>登录</h2>

            {/* 用户名输入 */}
            <div>
                <Label.Root htmlFor="username">
                    用户名
                </Label.Root>
                <input
                    id="username"
                    type="text"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    placeholder="请输入用户名"
                />
            </div>

            {/* 密码输入 */}
            <div>
                <Label.Root htmlFor="password">
                    密码
                </Label.Root>
                <div>
                    <input
                        id="password"
                        type={showPassword ? "text" : "password"}
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="请输入密码"
                    />
                    <button
                        type="button"
                        onClick={() => setShowPassword((v) => !v)}
                    >
                        {showPassword ? (
                            <EyeOpenIcon className="w-5 h-5" />
                        ) : (
                            <EyeClosedIcon className="w-5 h-5" />
                        )}
                    </button>{" "}
                </div>
            </div>

            {/* 用户协议 */}
            <div>
                <Checkbox.Root
                    checked={agreed}
                    onCheckedChange={(checked) => setAgreed(checked === true)}
                    id="terms"
                >
                    <Checkbox.Indicator>
                        <CheckIcon />
                    </Checkbox.Indicator>
                </Checkbox.Root>

                <label htmlFor="terms">
                    我已阅读并同意{" "}
                    <a href="/terms" >
                        《用户协议》
                    </a>
                </label>
            </div>

            {/* 登录按钮 */}
            <button
                type="submit"
                disabled={!username || !password || !agreed}
            >
                登录
            </button>
        </form>
    );
}
