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
      <div className={styles.valueContainer}>
        {loading ? (
          <Text>Loading values...</Text>
        ) : (
          values.map((value) => {
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
          })
        )}
      </div>

      <Text className={styles.hint}>
        Select 1-3 values ({selectedValueIds.length}/3 selected)
      </Text>
    </Field>
  );
};
