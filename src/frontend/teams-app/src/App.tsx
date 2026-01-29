import { useEffect, useState } from "react";
import { app, app as teamsApp } from "@microsoft/teams-js";
import {
  FluentProvider,
  webLightTheme,
  teamsLightTheme,
  teamsDarkTheme,
  teamsHighContrastTheme,
} from "@fluentui/react-components";
import type { Theme } from "@fluentui/react-components";
import { CardForm } from "./components/CardForm";
import { CardHistory } from "./components/CardHistory";
import "./App.css";

type TeamsTheme = "default" | "dark" | "contrast";

function App() {
  const [teamsTheme, setTeamsTheme] = useState<TeamsTheme>("default");
  const [initialized, setInitialized] = useState(false);
  const [currentView, setCurrentView] = useState<"create" | "history">(
    "create",
  );

  useEffect(() => {
    const initializeTeams = async () => {
      try {
        await teamsApp.initialize();

        // Get initial context
        const context = await app.getContext();
        setTeamsTheme(context.app.theme as TeamsTheme);

        // Check for deep link context
        if (context.page.subPageId === "history") {
          setCurrentView("history");
        } else {
          setCurrentView("create");
        }

        // Listen for theme changes
        teamsApp.registerOnThemeChangeHandler((theme: string) => {
          setTeamsTheme(theme as TeamsTheme);
        });

        // Listen for new theme changes (v2) if available
        // ...

        setInitialized(true);
      } catch (error) {
        console.error("Failed to initialize Teams SDK:", error);
        // Fall back to web theme if not in Teams
        setInitialized(true);
      }
    };

    initializeTeams();
  }, []);

  const getTheme = (): Theme => {
    // If running in Teams, use Teams themes
    if (initialized) {
      switch (teamsTheme) {
        case "dark":
          return teamsDarkTheme;
        case "contrast":
          return teamsHighContrastTheme;
        default:
          return teamsLightTheme;
      }
    }

    // Fallback to web themes for local development
    return webLightTheme;
  };

  if (!initialized) {
    return (
      <FluentProvider theme={webLightTheme}>
        <div style={{ padding: "20px", textAlign: "center" }}>
          Initializing...
        </div>
      </FluentProvider>
    );
  }

  return (
    <FluentProvider theme={getTheme()}>
      {currentView === "create" ? (
        <CardForm />
      ) : (
        <CardHistory onBack={() => setCurrentView("create")} />
      )}
    </FluentProvider>
  );
}

export default App;
