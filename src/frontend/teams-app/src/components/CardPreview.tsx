import React, { useState, useEffect, useRef } from "react";
import {
  Text,
  makeStyles,
  tokens,
  TabList,
  Tab,
} from "@fluentui/react-components";
import type { SelectTabData, SelectTabEvent } from "@fluentui/react-components";
import type { CompanyValue, Recipient } from "../types";
import { useDebounce } from "../hooks/useDebounce";
import "../styles/animations.css";

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
    backgroundColor: tokens.colorNeutralBackground1,
    maxWidth: "500px", // Mimic Teams card width constraint
    width: "100%",
    position: "relative",
    transformStyle: "preserve-3d",
    perspective: "1000px",
  },
  cardInner: {
    position: "relative",
    width: "100%",
    transformStyle: "preserve-3d",
    transition: "transform 600ms ease-in-out",
  },
  cardFlipped: {
    transform: "rotateY(180deg)",
  },
  cardFace: {
    backfaceVisibility: "hidden",
    WebkitBackfaceVisibility: "hidden",
  },
  cardFaceBack: {
    transform: "rotateY(180deg)",
  },
  // --- Content Card Styles ---
  contentHeader: {
    backgroundColor: tokens.colorNeutralBackground2,
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
    transition: "background-color 300ms ease-out",
  },
  message: {
    fontSize: tokens.fontSizeBase300,
    lineHeight: tokens.lineHeightBase300,
    color: tokens.colorNeutralForeground1,
    whiteSpace: "pre-wrap",
    transition: "background-color 300ms ease-out",
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
    transition: "all 300ms ease-out",
  },

  // --- Envelope Card Styles ---
  envelopeContainer: {
    backgroundColor: tokens.colorBrandBackground,
    padding: tokens.spacingVerticalL,
    color: tokens.colorNeutralForegroundOnBrand,
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
    color: tokens.colorNeutralForegroundOnBrand,
  },
  envelopeSubtitle: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForegroundOnBrand,
    opacity: 0.8,
  },
  actionSet: {
    padding: tokens.spacingVerticalM,
    backgroundColor: tokens.colorNeutralBackground1,
  },
  openButton: {
    backgroundColor: tokens.colorNeutralBackground1,
    color: tokens.colorBrandForeground1,
    border: `1px solid ${tokens.colorBrandStroke1}`,
    padding: "6px 20px",
    borderRadius: "4px",
    fontSize: tokens.fontSizeBase300,
    fontWeight: tokens.fontWeightSemibold,
    cursor: "not-allowed", // Preview only
    width: "100%",
  },
  // --- Animation Styles ---
  highlight: {
    backgroundColor: tokens.colorBrandBackgroundInvertedSelected,
    transition: "background-color 300ms ease-out",
  },
  fadeIn: {
    animation: "fadeIn 300ms ease-in-out",
  },
  slideTransition: {
    transition: "opacity 400ms ease-in-out, transform 400ms cubic-bezier(0.4, 0, 0.2, 1)",
  },
  slideOut: {
    opacity: 0,
    transform: "translateX(-20px)",
  },
  slideIn: {
    opacity: 1,
    transform: "translateX(0)",
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
  const [isFlipping, setIsFlipping] = useState(false);
  
  // Debounce the preview updates to optimize performance
  const debouncedRecipients = useDebounce(recipients, 150);
  const debouncedReason = useDebounce(recognitionReason, 150);
  const debouncedValues = useDebounce(selectedValues, 150);

  // Track previous values to detect changes for highlighting
  const prevRecipientsRef = useRef<Recipient[]>([]);
  const prevReasonRef = useRef<string>("");
  const prevValuesRef = useRef<CompanyValue[]>([]);
  
  // Track which sections changed for highlighting
  const [highlightRecipients, setHighlightRecipients] = useState(false);
  const [highlightReason, setHighlightReason] = useState(false);
  const [highlightValues, setHighlightValues] = useState(false);

  // Mock sender name for preview
  const senderName = "Youshan Li";

  // Detect changes and trigger highlight animations
  useEffect(() => {
    const recipientsChanged = 
      JSON.stringify(prevRecipientsRef.current) !== JSON.stringify(debouncedRecipients);
    const reasonChanged = prevReasonRef.current !== debouncedReason;
    const valuesChanged = 
      JSON.stringify(prevValuesRef.current) !== JSON.stringify(debouncedValues);

    if (recipientsChanged && prevRecipientsRef.current.length > 0) {
      setHighlightRecipients(true);
      setTimeout(() => setHighlightRecipients(false), 600);
    }

    if (reasonChanged && prevReasonRef.current !== "") {
      setHighlightReason(true);
      setTimeout(() => setHighlightReason(false), 600);
    }

    if (valuesChanged && prevValuesRef.current.length > 0) {
      setHighlightValues(true);
      setTimeout(() => setHighlightValues(false), 600);
    }

    // Update refs
    prevRecipientsRef.current = debouncedRecipients;
    prevReasonRef.current = debouncedReason;
    prevValuesRef.current = debouncedValues;
  }, [debouncedRecipients, debouncedReason, debouncedValues]);

  const onTabSelect = (_: SelectTabEvent, data: SelectTabData) => {
    const newMode = data.value as ViewMode;
    if (newMode !== viewMode) {
      setIsFlipping(true);
      // Trigger flip animation
      setTimeout(() => {
        setViewMode(newMode);
        setIsFlipping(false);
      }, 300); // Half of the flip duration for smooth transition
    }
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
        <div 
          className={`${styles.cardInner} ${isFlipping ? styles.cardFlipped : ''}`}
        >
          {viewMode === "recipient" ? (
            // --- Recipient (Envelope) View ---
            <div 
              className={`${styles.slideTransition} ${isFlipping ? styles.slideOut : styles.slideIn}`}
            >
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
            </div>
          ) : (
            // --- Public (Content) View ---
            <div 
              className={`${styles.slideTransition} ${isFlipping ? styles.slideOut : styles.slideIn}`}
            >
              <div className={styles.contentHeader}>
                <Text className={styles.contentTitle}>🎉 Recognition Card</Text>
              </div>
              <div className={styles.contentBody}>
                <div className={styles.factSet}>
                  <Text className={styles.factTitle}>To:</Text>
                  <Text 
                    className={`${styles.factValue} ${highlightRecipients ? styles.highlight : ''}`}
                  >
                    {debouncedRecipients.length > 0
                      ? debouncedRecipients.map((r) => r.name).join(", ")
                      : "Your Colleague"}
                  </Text>

                  <Text className={styles.factTitle}>From:</Text>
                  <Text className={styles.factValue}>{senderName}</Text>
                </div>

                <Text 
                  className={`${styles.message} ${highlightReason ? styles.highlight : ''}`}
                >
                  {debouncedReason ||
                    "Your appreciation message will appear here..."}
                </Text>

                {debouncedValues.length > 0 && (
                  <div className={styles.valuesContainer}>
                    {debouncedValues.map((v, index) => (
                      <span 
                        key={v.id} 
                        className={`${styles.valueTag} ${highlightValues ? styles.highlight : ''}`}
                        style={{
                          animationDelay: `${index * 50}ms`,
                        }}
                      >
                        {v.name}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
