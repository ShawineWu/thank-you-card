import React from "react";
import {
  ToggleButton,
  makeStyles,
  tokens,
  Text,
  Tooltip,
  Label,
} from "@fluentui/react-components";
import { CheckmarkCircleFilled } from "@fluentui/react-icons";
import { SkeletonLoader } from "./SkeletonLoader";
import type { CompanyValue } from "../types";
import { getValueColor } from "../constants/valueColors";

const useStyles = makeStyles({
  container: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalM,
    marginBottom: tokens.spacingVerticalL,
  },
  section: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalS,
  },
  sectionLabel: {
    fontWeight: tokens.fontWeightSemibold,
    fontSize: tokens.fontSizeBase300,
  },
  valueContainer: {
    display: "flex",
    flexWrap: "wrap",
    gap: tokens.spacingHorizontalS,
  },
  valueChip: {
    minWidth: "120px",
  },
  hint: {
    marginTop: tokens.spacingVerticalXS,
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground3,
  },
  errorText: {
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorPaletteRedForeground1,
  },
  loadingContainer: {
    display: "flex",
    flexWrap: "wrap",
    gap: tokens.spacingHorizontalS,
    opacity: 0,
    animation: "fadeIn 300ms ease-in-out forwards",
  },
  valuesLoaded: {
    opacity: 0,
    animation: "fadeIn 400ms ease-in-out forwards",
  },
});

interface ValueSelectorProps {
  values: CompanyValue[];
  selectedValueIds: string[];
  onChange: (valueIds: string[]) => void;
  error?: string;
  loading?: boolean;
}

export const ValueSelector: React.FC<ValueSelectorProps> = ({
  values,
  selectedValueIds,
  onChange,
  error,
  loading,
}) => {
  const styles = useStyles();

  const handleToggle = (valueId: string) => {
    const isSelected = selectedValueIds.includes(valueId);

    if (isSelected) {
      onChange(selectedValueIds.filter((id) => id !== valueId));
    } else {
      if (selectedValueIds.length < 3) {
        onChange([...selectedValueIds, valueId]);
      }
    }
  };

  // Separate values by type
  const credos = values.filter((v) => v.type === "Credo");
  const coreValues = values.filter((v) => v.type === "Value");

  const renderValueButtons = (items: CompanyValue[]) => (
    <div className={`${styles.valueContainer} ${styles.valuesLoaded}`}>
      {items.map((value) => {
        const isSelected = selectedValueIds.includes(value.id);
        const isDisabled = !isSelected && selectedValueIds.length >= 3;
        const colorConfig = getValueColor(value.name);

        return (
          <Tooltip
            key={value.id}
            content={value.description}
            relationship="description"
            positioning="above"
          >
            <ToggleButton
              checked={isSelected}
              onClick={() => handleToggle(value.id)}
              disabled={isDisabled}
              className={styles.valueChip}
              icon={isSelected ? <CheckmarkCircleFilled style={{ color: colorConfig.textCss }} /> : undefined}
              appearance="subtle"
              style={{
                backgroundColor: colorConfig.bgCss,
                color: colorConfig.textCss,
                borderColor: isSelected ? colorConfig.textCss : colorConfig.borderCss,
                borderWidth: isSelected ? "2px" : "1px",
                boxShadow: isSelected ? `0 2px 8px ${colorConfig.borderCss}` : undefined,
                opacity: isDisabled ? 0.5 : 1,
              }}
            >
              {value.name}
            </ToggleButton>
          </Tooltip>
        );
      })}
    </div>
  );

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.section}>
          <Label className={styles.sectionLabel}>Credo</Label>
          <div className={styles.loadingContainer}>
            <SkeletonLoader variant="rectangular" width="120px" height="32px" />
            <SkeletonLoader variant="rectangular" width="120px" height="32px" />
            <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          </div>
        </div>
        <div className={styles.section}>
          <Label className={styles.sectionLabel}>Value</Label>
          <div className={styles.loadingContainer}>
            <SkeletonLoader variant="rectangular" width="120px" height="32px" />
            <SkeletonLoader variant="rectangular" width="120px" height="32px" />
            <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Credo Section */}
      {credos.length > 0 && (
        <div className={styles.section}>
          <Label className={styles.sectionLabel} required>
            Credo
          </Label>
          {renderValueButtons(credos)}
        </div>
      )}

      {/* Value Section */}
      {coreValues.length > 0 && (
        <div className={styles.section}>
          <Label className={styles.sectionLabel}>Value</Label>
          {renderValueButtons(coreValues)}
        </div>
      )}

      {/* Hint and Error */}
      <div>
        <Text className={styles.hint}>
          Select 1-3 values ({selectedValueIds.length}/3 selected)
        </Text>
        {error && <Text className={styles.errorText}>{error}</Text>}
      </div>
    </div>
  );
};
