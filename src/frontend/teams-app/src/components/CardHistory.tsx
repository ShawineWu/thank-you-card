import React from "react";
import { makeStyles, tokens, Text, Button } from "@fluentui/react-components";
import { ArrowLeftRegular } from "@fluentui/react-icons";

const useStyles = makeStyles({
  container: {
    padding: tokens.spacingHorizontalL,
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalM,
  },
  header: {
    display: "flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalM,
    marginBottom: tokens.spacingVerticalL,
  },
  emptyState: {
    textAlign: "center",
    padding: tokens.spacingVerticalXXL,
    color: tokens.colorNeutralForeground3,
    backgroundColor: tokens.colorNeutralBackground2,
    borderRadius: tokens.borderRadiusMedium,
  },
});

interface CardHistoryProps {
  onBack: () => void;
}

export const CardHistory: React.FC<CardHistoryProps> = ({ onBack }) => {
  const styles = useStyles();

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Button
          icon={<ArrowLeftRegular />}
          appearance="subtle"
          onClick={onBack}
          aria-label="Back to create"
        />
        <Text size={600} weight="semibold">
          Praise History
        </Text>
      </div>

      <div className={styles.emptyState}>
        <Text>You haven't received any praise yet... but stay awesome! 🌟</Text>
        {/* TODO: Integrate with backend to fetch actual history */}
      </div>
    </div>
  );
};
