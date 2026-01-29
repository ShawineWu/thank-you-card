import React from "react";
import { NavLink } from "react-router-dom";
import {
  LayoutDashboard,
  Send,
  Inbox,
  BarChart3,
  Settings,
  Heart,
} from "lucide-react";
import { cn } from "@/utils/cn";

export const Sidebar: React.FC = () => {
  const navItems = [
    {
      icon: <LayoutDashboard size={20} />,
      label: "Overview",
      path: "/dashboard",
    },
    { icon: <Inbox size={20} />, label: "Received Cards", path: "/received" },
    { icon: <Send size={20} />, label: "Sent Cards", path: "/sent" },
    { icon: <BarChart3 size={20} />, label: "Statistics", path: "/stats" },
    {
      icon: <BarChart3 size={20} />,
      label: "HR Analytics",
      path: "/analytics",
    },
  ];

  return (
    <aside className="fixed left-0 top-0 h-screen w-[260px] bg-card border-r border-border flex flex-col z-[101]">
      <div className="h-[70px] px-6 flex items-center gap-3">
        <Heart className="text-primary h-6 w-6" fill="currentColor" />
        <span className="text-xl font-bold tracking-tight">ThankYou</span>
      </div>

      <nav className="flex-1 p-4 flex flex-col gap-8 overflow-y-auto">
        <div className="space-y-2">
          <p className="px-3 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
            Main Menu
          </p>
          <ul className="space-y-1">
            {navItems.map((item) => (
              <li key={item.path}>
                <NavLink
                  to={item.path}
                  className={({ isActive }) =>
                    cn(
                      "flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all font-medium",
                      isActive
                        ? "bg-primary text-primary-foreground shadow-lg shadow-primary/20"
                        : "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
                    )
                  }
                >
                  {item.icon}
                  <span>{item.label}</span>
                </NavLink>
              </li>
            ))}
          </ul>
        </div>

        <div className="mt-auto pt-4 border-t border-border">
          <p className="px-3 text-xs font-semibold text-muted-foreground uppercase tracking-wider mb-2">
            Others
          </p>
          <NavLink
            to="/settings"
            className={({ isActive }) =>
              cn(
                "flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all font-medium",
                isActive
                  ? "bg-primary text-primary-foreground"
                  : "text-muted-foreground hover:bg-accent hover:text-accent-foreground",
              )
            }
          >
            <Settings size={20} />
            <span>Settings</span>
          </NavLink>
        </div>
      </nav>
    </aside>
  );
};
