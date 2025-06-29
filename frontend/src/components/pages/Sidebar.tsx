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
export const Sidebar = () => {
    <Sidebar>
        <SidebarHeader>AnnotateX</SidebarHeader>
        <SidebarContent>
            <SidebarMenu>
                <SidebarGroup>
                    <SidebarGroupLabel>Projects</SidebarGroupLabel>
                    <SidebarGroupContent>
                        <SidebarMenuItem>
                            <SidebarMenuButton asChild>
                                <a href="/project">
                                    <PanelsTopLeft />
                                    <span>Projects</span>
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarMenu>
        </SidebarContent>
    </Sidebar>;
};
