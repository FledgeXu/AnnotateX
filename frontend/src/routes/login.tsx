import { createFileRoute } from "@tanstack/react-router";
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
} from "@radix-ui/themes";
import { EyeOpenIcon } from "@radix-ui/react-icons";

export const Route = createFileRoute("/login")({
    component: Login,
});

function Login() {
    return (
        <Flex align="center" justify="center" height="100vh">
            <Card size="3" className="w-full max-w-sm">
                <form>
                    <Flex gap={"3"} direction={"column"}>
                        <Heading>Login</Heading>
                        <Box>
                            <Text weight={"medium"}>Username</Text>
                            <TextField.Root type="text">
                                <TextField.Slot />
                            </TextField.Root>
                        </Box>
                        <Box>
                            <Text weight={"medium"}>Password</Text>
                            <TextField.Root type="password">
                                <TextField.Slot />
                                <TextField.Slot>
                                    <IconButton variant="ghost">
                                        <EyeOpenIcon />
                                    </IconButton>
                                </TextField.Slot>
                            </TextField.Root>
                        </Box>
                        <Flex align="center" gap="2">
                            <Checkbox />
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
