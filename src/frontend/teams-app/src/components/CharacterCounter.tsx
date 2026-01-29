import React from "react";
import { makeStyles, tokens, Text } from "@fluentui/react-components";

export interface CharacterCounterProps {
  current: number;
  max: number;
  warningThreshold?: number; // 警告阈值（默认80%）
  errorThreshold?: number; // 错误阈值（默认100%）
}

export type CounterState = "normal" | "warning" | "error";

const useStyles = makeStyles({
  counter: {
    display: "inline-flex",
    alignItems: "center",
    fontSize: tokens.fontSizeBase200,
    fontWeight: tokens.fontWeightRegular,
    transition: "color 0.2s ease-in-out",
  },
  normal: {
    color: tokens.colorNeutralForeground3,
  },
  warning: {
    color: tokens.colorPaletteYellowForeground1,
  },
  error: {
    color: tokens.colorPaletteRedForeground1,
  },
});

export const CharacterCounter: React.FC<CharacterCounterProps> = ({
  current,
  max,
  warningThreshold = 0.8,
  errorThreshold = 1.0,
}) => {
  const styles = useStyles();

  // Calculate the state based on thresholds
  const getCounterState = (): CounterState => {
    const percentage = current / max;

    if (percentage >= errorThreshold) {
      return "error";
    } else if (percentage >= warningThreshold) {
      return "warning";
    }
    return "normal";
  };

  const state = getCounterState();
  const remaining = max - current;

  // Determine the CSS class based on state
  const stateClass =
    state === "error"
      ? styles.error
      : state === "warning"
        ? styles.warning
        : styles.normal;

  return (
    <Text className={`${styles.counter} ${stateClass}`}>
      {remaining >= 0 ? `${remaining} characters remaining` : `${Math.abs(remaining)} characters over limit`}
    </Text>
  );
};
