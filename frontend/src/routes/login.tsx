import { createFileRoute } from "@tanstack/react-router";
import { CardHeader, CardTitle, Card, CardContent } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { LoginForm } from "@/components/pages/LoginForm";

export const Route = createFileRoute("/login")({
  component: Login,
});

function Login() {
  return (
    <div className="flex items-center justify-center h-screen">
      <Card className="w-full max-w-sm">
        <CardHeader>
          <CardTitle>
            <Label className="text-xl font-bold">Login</Label>
          </CardTitle>
        </CardHeader>
        <CardContent>
          <LoginForm />
        </CardContent>
      </Card>
    </div>
  );
}
