import React, { useState } from "react";
import {
  Text,
  makeStyles,
  tokens,
  TabList,
  Tab,
} from "@fluentui/react-components";
import type { SelectTabData, SelectTabEvent } from "@fluentui/react-components";
import type { CompanyValue, Recipient } from "../types";

const useStyles = makeStyles({
  container: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalM,
    marginTop: tokens.spacingVerticalL,
  },
  card: {
    borderRadius: "4px", // Adaptive cards are usually sharper
    boxShadow: tokens.shadow8,
    border: `1px solid ${tokens.colorNeutralStroke2}`,
    overflow: "hidden",
    backgroundColor: "white",
    maxWidth: "500px", // Mimic Teams card width constraint
    width: "100%",
  },
  // --- Content Card Styles ---
  contentHeader: {
    backgroundColor: "#F0F0F0", // "Emphasis" style
    padding: tokens.spacingVerticalM,
    borderBottom: `1px solid ${tokens.colorNeutralStroke2}`,
  },
  contentTitle: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightBold,
    color: tokens.colorBrandForeground1,
  },
  contentBody: {
    padding: tokens.spacingVerticalM,
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalM,
  },
  factSet: {
    display: "grid",
    gridTemplateColumns: "auto 1fr",
    columnGap: tokens.spacingHorizontalL,
    rowGap: tokens.spacingVerticalS,
  },
  factTitle: {
    fontWeight: tokens.fontWeightSemibold,
    color: tokens.colorNeutralForeground1,
  },
  factValue: {
    color: tokens.colorNeutralForeground1,
  },
  message: {
    fontSize: tokens.fontSizeBase300,
    lineHeight: tokens.lineHeightBase300,
    color: tokens.colorNeutralForeground1,
    whiteSpace: "pre-wrap",
  },
  valuesContainer: {
    display: "flex",
    gap: tokens.spacingHorizontalS,
    flexWrap: "wrap",
    marginTop: tokens.spacingVerticalS,
  },
  valueTag: {
    fontSize: "12px",
    padding: "2px 8px",
    backgroundColor: tokens.colorNeutralBackground2,
    borderRadius: "4px",
    color: tokens.colorNeutralForeground2,
  },

  // --- Envelope Card Styles ---
  envelopeContainer: {
    backgroundColor: "#6264A7", // "Accent" style (using Teams blurple)
    padding: tokens.spacingVerticalL,
    color: "white",
    display: "flex",
    flexDirection: "row", // ColumnSet
    alignItems: "center",
    gap: tokens.spacingHorizontalM,
  },
  iconColumn: {
    fontSize: "48px",
  },
  textColumn: {
    display: "flex",
    flexDirection: "column",
    justifyContent: "center",
    flex: 1,
  },
  envelopeTitle: {
    fontSize: tokens.fontSizeBase400,
    fontWeight: tokens.fontWeightBold,
    color: "white",
  },
  envelopeSubtitle: {
    fontSize: tokens.fontSizeBase200,
    color: "rgba(255, 255, 255, 0.8)", // "IsSubtle" on accent
  },
  actionSet: {
    padding: tokens.spacingVerticalM,
    backgroundColor: "white", // Actions are usually separate or in the container
  },
  openButton: {
    backgroundColor: "white",
    color: "#6264A7",
    border: "1px solid #6264A7",
    padding: "6px 20px",
    borderRadius: "4px",
    fontSize: tokens.fontSizeBase300,
    fontWeight: tokens.fontWeightSemibold,
    cursor: "not-allowed", // Preview only
    width: "100%",
  },
});

interface CardPreviewProps {
  recipients: Recipient[];
  recognitionReason: string;
  selectedValues: CompanyValue[];
  senderName?: string;
}

type ViewMode = "recipient" | "public";

export const CardPreview: React.FC<CardPreviewProps> = ({
  recipients,
  recognitionReason,
  selectedValues,
}) => {
  const styles = useStyles();
  const [viewMode, setViewMode] = useState<ViewMode>("public");

  // Mock sender name for preview
  const senderName = "Youshan Li";

  const onTabSelect = (_: SelectTabEvent, data: SelectTabData) => {
    setViewMode(data.value as ViewMode);
  };

  return (
    <div className={styles.container}>
      <div
        style={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <Text weight="semibold" size={400}>
          Preview
        </Text>
        <TabList
          selectedValue={viewMode}
          onTabSelect={onTabSelect}
          appearance="subtle"
        >
          <Tab value="public">Public View</Tab>
          <Tab value="recipient">Recipient View</Tab>
        </TabList>
      </div>

      <div className={styles.card}>
        {viewMode === "recipient" ? (
          // --- Recipient (Envelope) View ---
          <>
            <div className={styles.envelopeContainer}>
              <div className={styles.iconColumn}>✉️</div>
              <div className={styles.textColumn}>
                <Text className={styles.envelopeTitle}>
                  You've received a Thank You Card!
                </Text>
                <Text className={styles.envelopeSubtitle}>
                  From {senderName}
                </Text>
              </div>
            </div>
            <div className={styles.actionSet}>
              <button className={styles.openButton}>Open Card</button>
            </div>
          </>
        ) : (
          // --- Public (Content) View ---
          <>
            <div className={styles.contentHeader}>
              <Text className={styles.contentTitle}>🎉 Recognition Card</Text>
            </div>
            <div className={styles.contentBody}>
              <div className={styles.factSet}>
                <Text className={styles.factTitle}>To:</Text>
                <Text className={styles.factValue}>
                  {recipients.length > 0
                    ? recipients.map((r) => r.name).join(", ")
                    : "Your Colleague"}
                </Text>

                <Text className={styles.factTitle}>From:</Text>
                <Text className={styles.factValue}>{senderName}</Text>
              </div>

              <Text className={styles.message}>
                {recognitionReason ||
                  "Your appreciation message will appear here..."}
              </Text>

              {selectedValues.length > 0 && (
                <div className={styles.valuesContainer}>
                  {selectedValues.map((v) => (
                    <span key={v.id} className={styles.valueTag}>
                      {v.name}
                    </span>
                  ))}
                </div>
              )}
            </div>
          </>
        )}
      </div>
    </div>
  );
};
