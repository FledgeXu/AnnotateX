import { Link } from "@tanstack/react-router";
import { PanelsTopLeft } from "lucide-react";
import {
    Sidebar,
    SidebarContent,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarGroup,
    SidebarGroupLabel,
    SidebarGroupContent,
} from "@/components/ui/sidebar";

export const HomePageSidebar = () => {
    return (
        <Sidebar>
            <SidebarHeader>AnnotateX</SidebarHeader>
            <SidebarContent>
                <SidebarMenu>
                    <SidebarGroup>
                        <SidebarGroupLabel>Projects</SidebarGroupLabel>
                        <SidebarGroupContent>
                            <SidebarMenuItem>
                                <SidebarMenuButton asChild>
                                    <Link to="/project">
                                        <PanelsTopLeft />
                                        <span>Projects</span>
                                    </Link>
                                </SidebarMenuButton>
                            </SidebarMenuItem>
                        </SidebarGroupContent>
                    </SidebarGroup>
                </SidebarMenu>
            </SidebarContent>
        </Sidebar>
    );
};
