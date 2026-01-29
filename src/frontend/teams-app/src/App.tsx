import { useEffect, useState } from "react";
import { app, app as teamsApp } from "@microsoft/teams-js";
import {
  FluentProvider,
  webLightTheme,
  teamsLightTheme,
  teamsDarkTheme,
  teamsHighContrastTheme,
  makeStyles,
  tokens,
} from "@fluentui/react-components";
import type { Theme } from "@fluentui/react-components";
import { CardForm } from "./components/CardForm";
import { CardHistory } from "./components/CardHistory";
import { SkeletonLoader } from "./components/SkeletonLoader";
import "./App.css";

type TeamsTheme = "default" | "dark" | "contrast";

const useStyles = makeStyles({
  skeletonContainer: {
    padding: tokens.spacingVerticalXL,
    maxWidth: "800px",
    margin: "0 auto",
    opacity: 0,
    animation: "fadeIn 300ms ease-in-out forwards",
  },
  skeletonHeader: {
    marginBottom: tokens.spacingVerticalL,
  },
  skeletonForm: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalL,
  },
  skeletonField: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalS,
  },
  skeletonActions: {
    display: "flex",
    gap: tokens.spacingHorizontalM,
    marginTop: tokens.spacingVerticalL,
  },
  contentContainer: {
    opacity: 0,
    animation: "fadeIn 400ms ease-in-out forwards",
  },
});

function App() {
  const styles = useStyles();
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
        <div className={styles.skeletonContainer}>
          {/* Header skeleton */}
          <div className={styles.skeletonHeader}>
            <SkeletonLoader variant="text" width="60%" height="32px" />
            <div style={{ marginTop: tokens.spacingVerticalS }}>
              <SkeletonLoader variant="text" width="80%" height="16px" />
            </div>
          </div>

          {/* Form skeleton */}
          <div className={styles.skeletonForm}>
            {/* Recipient selector skeleton */}
            <div className={styles.skeletonField}>
              <SkeletonLoader variant="text" width="150px" height="16px" />
              <SkeletonLoader variant="rectangular" height="40px" />
            </div>

            {/* Recognition reason skeleton */}
            <div className={styles.skeletonField}>
              <SkeletonLoader variant="text" width="180px" height="16px" />
              <SkeletonLoader variant="rectangular" height="120px" />
              <SkeletonLoader variant="text" width="100px" height="14px" />
            </div>

            {/* Value selector skeleton */}
            <div className={styles.skeletonField}>
              <SkeletonLoader variant="text" width="140px" height="16px" />
              <div style={{ 
                display: "flex", 
                flexWrap: "wrap", 
                gap: tokens.spacingHorizontalS,
                marginTop: tokens.spacingVerticalS 
              }}>
                <SkeletonLoader variant="rectangular" width="120px" height="32px" />
                <SkeletonLoader variant="rectangular" width="120px" height="32px" />
                <SkeletonLoader variant="rectangular" width="120px" height="32px" />
                <SkeletonLoader variant="rectangular" width="120px" height="32px" />
              </div>
            </div>

            {/* Preview skeleton */}
            <div className={styles.skeletonField}>
              <SkeletonLoader variant="text" width="120px" height="16px" />
              <SkeletonLoader variant="rectangular" height="200px" />
            </div>

            {/* Actions skeleton */}
            <div className={styles.skeletonActions}>
              <SkeletonLoader variant="rectangular" width="180px" height="40px" />
              <SkeletonLoader variant="rectangular" width="100px" height="40px" />
            </div>
          </div>
        </div>
      </FluentProvider>
    );
  }

  return (
    <FluentProvider theme={getTheme()}>
      <div className={styles.contentContainer}>
        {currentView === "create" ? (
          <CardForm />
        ) : (
          <CardHistory onBack={() => setCurrentView("create")} />
        )}
      </div>
    </FluentProvider>
  );
}

export default App;
