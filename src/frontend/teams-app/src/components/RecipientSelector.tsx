import React, { useState, useEffect } from "react";
import {
  Combobox,
  Option,
  Tag,
  TagGroup,
  Field,
  makeStyles,
  tokens,
  Spinner,
} from "@fluentui/react-components";
import { PersonRegular } from "@fluentui/react-icons";
import { api } from "../services/api";
import type { Employee, Recipient } from "../types";

const useStyles = makeStyles({
  field: {
    marginBottom: tokens.spacingVerticalL,
  },
  tagGroup: {
    marginTop: tokens.spacingVerticalS,
  },
});

interface RecipientSelectorProps {
  selectedRecipients: Recipient[];
  onChange: (recipients: Recipient[]) => void;
  error?: string;
}

export const RecipientSelector: React.FC<RecipientSelectorProps> = ({
  selectedRecipients,
  onChange,
  error,
}) => {
  const styles = useStyles();
  const [searchValue, setSearchValue] = useState("");
  const [searchResults, setSearchResults] = useState<Employee[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  // Debounced search
  useEffect(() => {
    const search = async () => {
      if (searchValue.trim().length < 2) {
        setSearchResults([]);
        return;
      }

      setIsLoading(true);
      try {
        const results = await api.searchEmployees(searchValue.trim());
        setSearchResults(results);
      } catch (err) {
        console.error("Failed to search employees", err);
      } finally {
        setIsLoading(false);
      }
    };

    const timeoutId = setTimeout(search, 300);
    return () => clearTimeout(timeoutId);
  }, [searchValue]);

  const handleSelect = (_event: any, data: any) => {
    const employeeId = data.optionValue;
    const employee = searchResults.find((emp) => emp.id === employeeId);

    if (employee && !selectedRecipients.find((r) => r.id === employeeId)) {
      onChange([
        ...selectedRecipients,
        { id: employee.id, name: employee.name },
      ]);
    }
    setSearchValue("");
  };

  const handleRemove = (employeeId: string) => {
    onChange(selectedRecipients.filter((r) => r.id !== employeeId));
  };

  return (
    <Field
      label={
        <>
          Recipients <span style={{ color: "red" }}>*</span>
        </>
      }
      validationMessage={error}
      validationState={error ? "error" : "none"}
      className={styles.field}
    >
      <Combobox
        placeholder="Search for colleagues..."
        value={searchValue}
        onInput={(e) => setSearchValue(e.currentTarget.value)}
        onOptionSelect={handleSelect}
      >
        {isLoading && <Spinner size="tiny" label="Searching..." />}
        {searchResults.map((emp) => (
          <Option key={emp.id} text={emp.name} value={emp.id}>
            <PersonRegular style={{ marginRight: "8px" }} />
            <div>
              <div>{emp.name}</div>
              <div
                style={{
                  fontSize: "12px",
                  color: tokens.colorNeutralForeground3,
                }}
              >
                {emp.email} {emp.department ? `· ${emp.department}` : ""}
              </div>
            </div>
          </Option>
        ))}
      </Combobox>

      {selectedRecipients.length > 0 && (
        <TagGroup
          className={styles.tagGroup}
          onDismiss={(_e, { value }) => {
            if (value) handleRemove(value);
          }}
        >
          {selectedRecipients.map((emp) => (
            <Tag
              key={emp.id}
              value={emp.id}
              dismissible
              dismissIcon={{ "aria-label": "remove" }}
            >
              {emp.name}
            </Tag>
          ))}
        </TagGroup>
      )}
    </Field>
  );
};
