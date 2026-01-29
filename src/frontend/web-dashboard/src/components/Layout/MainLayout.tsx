import React from "react";
import { Outlet } from "react-router-dom";
import { Sidebar } from "../Sidebar/Sidebar";
import { Header } from "../Header/Header";

export const MainLayout: React.FC = () => {
  return (
    <div className="flex min-h-screen bg-background text-foreground">
      <Sidebar />
      <div className="flex-1 ml-[260px] flex flex-col">
        <Header />
        <main className="flex-1 p-8 overflow-y-auto bg-[radial-gradient(at_0%_0%,rgba(99,102,241,0.03)_0px,transparent_50%),radial-gradient(at_100%_0%,rgba(236,72,153,0.03)_0px,transparent_50%)]">
          <div className="max-w-[1200px] mx-auto">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};
