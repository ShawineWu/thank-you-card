import React from "react";
import { useAuthStore } from "@/store/authStore";
import { LogOut, Bell } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { initConfigProcess } from "@/app";

export const Header: React.FC = () => {
  const { user } = useAuthStore();
  const userName =
    user?.profile.name || user?.profile.preferred_username || "User";

  const getInitials = (name: string) => {
    return (
      name
        ?.split(" ")
        .map((n) => n[0])
        .join("")
        .toUpperCase() || "U"
    );
  };

  const handleLogout = async () => {
    try {
      const auth = await initConfigProcess;
      await auth.logout();
    } catch (error) {
      console.error("[Header] Logout failed:", error);
    }
  };

  return (
    <header className="sticky top-0 z-50 h-[70px] w-full border-b border-border bg-background/80 backdrop-blur-md px-8 flex items-center justify-between">
      <div className="flex items-center gap-4">
        <h2 className="text-lg font-semibold bg-gradient-to-r from-foreground to-muted-foreground bg-clip-text text-transparent">
          Dashboard
        </h2>
      </div>

      <div className="flex items-center gap-6">
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-foreground"
        >
          <Bell size={20} />
        </Button>

        <div className="flex items-center gap-3 pr-4 border-r border-border">
          <div className="text-right hidden sm:block">
            <p className="text-sm font-semibold leading-none">{userName}</p>
            <p className="text-[10px] text-muted-foreground font-medium uppercase mt-1 tracking-wider">
              Employee
            </p>
          </div>
          <Avatar className="h-9 w-9 border-2 border-primary/20">
            <AvatarFallback className="bg-primary/10 text-primary font-bold">
              {getInitials(userName)}
            </AvatarFallback>
          </Avatar>
        </div>

        <Button
          variant="destructive"
          size="sm"
          onClick={handleLogout}
          className="gap-2 rounded-full px-4 font-medium h-9"
        >
          <LogOut size={16} />
          <span className="hidden sm:inline">Logout</span>
        </Button>
      </div>
    </header>
  );
};
