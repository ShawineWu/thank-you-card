import React from "react";
import {
  Field,
  ToggleButton,
  makeStyles,
  tokens,
  Text,
  Tooltip,
} from "@fluentui/react-components";
import { CheckmarkCircleFilled } from "@fluentui/react-icons";
import { SkeletonLoader } from "./SkeletonLoader";
import type { CompanyValue } from "../types";

const useStyles = makeStyles({
  field: {
    marginBottom: tokens.spacingVerticalL,
  },
  valueContainer: {
    display: "flex",
    flexWrap: "wrap",
    gap: tokens.spacingHorizontalS,
    marginTop: tokens.spacingVerticalS,
  },
  valueChip: {
    minWidth: "120px",
  },
  hint: {
    marginTop: tokens.spacingVerticalXS,
    fontSize: tokens.fontSizeBase200,
    color: tokens.colorNeutralForeground3,
  },
  loadingContainer: {
    display: "flex",
    flexWrap: "wrap",
    gap: tokens.spacingHorizontalS,
    marginTop: tokens.spacingVerticalS,
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

  return (
    <Field
      label="Company Values"
      required
      validationMessage={error}
      validationState={error ? "error" : "none"}
      className={styles.field}
    >
      {loading ? (
        <div className={styles.loadingContainer}>
          <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          <SkeletonLoader variant="rectangular" width="120px" height="32px" />
          <SkeletonLoader variant="rectangular" width="120px" height="32px" />
        </div>
      ) : (
        <div className={`${styles.valueContainer} ${styles.valuesLoaded}`}>
          {values.map((value) => {
            const isSelected = selectedValueIds.includes(value.id);
            const isDisabled = !isSelected && selectedValueIds.length >= 3;

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
                  icon={isSelected ? <CheckmarkCircleFilled /> : undefined}
                  appearance={isSelected ? "primary" : "subtle"}
                >
                  {value.name}
                </ToggleButton>
              </Tooltip>
            );
          })}
        </div>
      )}

      <Text className={styles.hint}>
        Select 1-3 values ({selectedValueIds.length}/3 selected)
      </Text>
    </Field>
  );
};
