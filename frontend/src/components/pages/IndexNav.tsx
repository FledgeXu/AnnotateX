import { useNavigate } from "@tanstack/react-router";
import { useStoreActions } from "easy-peasy";
import { CircleUserRound, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import type { StoreModel } from "@/store/types";

export const IndexNav = () => {
  const authLogout = useStoreActions<StoreModel>(
    (actions) => actions.auth.logout,
  );
  const navigate = useNavigate();
  const logout = async () => {
    authLogout();
    navigate({ to: "/login" });
  };
  return (
    <nav className="flex items-center justify-between px-4 py-3">
      <SidebarTrigger className="w-5 h-5" />
      <div className="flex items-center gap-4 h-6">
        <Button
          variant={"secondary"}
          className="bg-gray-100 border h-6 w-40 text-gray-500"
        >
          <Search />
          <span>Search Project</span>
        </Button>
        <Separator orientation="vertical" />
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <CircleUserRound className="w-5 h-5" />
          </DropdownMenuTrigger>
          <DropdownMenuContent>
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>Profile</DropdownMenuItem>
            <DropdownMenuItem onClick={logout}>Logout</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </nav>
  );
};
