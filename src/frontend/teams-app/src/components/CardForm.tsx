import React, { useState, useEffect } from "react";
import {
  Button,
  Textarea,
  Field,
  makeStyles,
  tokens,
  Text,
  MessageBar,
  MessageBarBody,
  MessageBarTitle,
} from "@fluentui/react-components";
import { SendRegular } from "@fluentui/react-icons";
import * as microsoftTeams from "@microsoft/teams-js";
import { RecipientSelector } from "./RecipientSelector";
import { ValueSelector } from "./ValueSelector";
import { CardPreview } from "./CardPreview";
import { api } from "../services/api";
import type { CompanyValue, CardFormData, Recipient } from "../types";

const useStyles = makeStyles({
  container: {
    padding: tokens.spacingVerticalXL,
    maxWidth: "800px",
    margin: "0 auto",
  },
  header: {
    marginBottom: tokens.spacingVerticalL,
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: tokens.spacingVerticalL,
  },
  textareaField: {
    marginBottom: tokens.spacingVerticalL,
  },
  actions: {
    display: "flex",
    gap: tokens.spacingHorizontalM,
    marginTop: tokens.spacingVerticalL,
  },
  message: {
    marginBottom: tokens.spacingVerticalM,
  },
});

export const CardForm: React.FC = () => {
  const styles = useStyles();
  const [formData, setFormData] = useState<CardFormData>({
    recipients: [],
    recognitionReason: "",
    valueIds: [],
  });
  const [errors, setErrors] = useState<
    Partial<Record<keyof CardFormData, string>>
  >({});
  const [loading, setLoading] = useState(false);
  const [successMessage, setSuccessMessage] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [companyValues, setCompanyValues] = useState<CompanyValue[]>([]);
  const [loadingValues, setLoadingValues] = useState(true);
  const [isInTeamsTask, setIsInTeamsTask] = useState(false);

  // Initialize Teams SDK
  useEffect(() => {
    try {
      microsoftTeams.app.initialize().then(() => {
        microsoftTeams.app.getContext().then((context) => {
          // If we are in a task module or compose extension
          if (
            context.page.frameContext === "task" ||
            context.page.frameContext === "content"
          ) {
            setIsInTeamsTask(true);
          }
        });
      });
    } catch {
      console.log("Not running in Teams context");
    }
  }, []);

  // Fetch company values from API
  useEffect(() => {
    const fetchValues = async () => {
      try {
        setLoadingValues(true);
        const values = await api.getCompanyValues();
        setCompanyValues(values);
      } catch (error: any) {
        setErrorMessage(error.message || "Failed to load company values");
      } finally {
        setLoadingValues(false);
      }
    };
    fetchValues();
  }, []);

  const validate = (): boolean => {
    const newErrors: Partial<Record<keyof CardFormData, string>> = {};

    if (formData.recipients.length === 0) {
      newErrors.recipients = "Please select at least one recipient";
    }

    if (!formData.recognitionReason.trim()) {
      newErrors.recognitionReason = "Please provide a recognition reason";
    } else if (formData.recognitionReason.trim().length < 10) {
      newErrors.recognitionReason = "Reason must be at least 10 characters";
    } else if (formData.recognitionReason.length > 500) {
      newErrors.recognitionReason = "Reason must not exceed 500 characters";
    }

    if (formData.valueIds.length === 0) {
      newErrors.valueIds = "Please select at least one company value";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage("");
    setErrorMessage("");

    if (!validate()) {
      return;
    }

    try {
      setLoading(true);
      const cardRequest = {
        recipients: formData.recipients,
        recognitionReason: formData.recognitionReason.trim(),
        valueIds: formData.valueIds,
      };

      await api.createCard(cardRequest);

      setSuccessMessage("Recognition card sent successfully! 🎉");

      // If running in a Teams Task Module (Message Extension), submit the task
      // This will close the popup and pass the data back to our Bot backend
      if (isInTeamsTask) {
        setLoading(true); // Keep loading while closing
        setTimeout(() => {
          microsoftTeams.tasks.submitTask({
            data: {
              ...cardRequest,
              senderName: "Youshan Li", // Pass sender name for the Adaptive Card
            },
          });
        }, 1000);
        return;
      }

      // Reset form if not in a task
      setFormData({
        recipients: [],
        recognitionReason: "",
        valueIds: [],
      });
      setErrors({});
    } catch (error: any) {
      setErrorMessage(error.message || "Failed to send recognition card");
    } finally {
      if (!isInTeamsTask) {
        setLoading(false);
      }
    }
  };

  const handleReset = () => {
    setFormData({
      recipients: [],
      recognitionReason: "",
      valueIds: [],
    });
    setErrors({});
    setSuccessMessage("");
    setErrorMessage("");
  };

  const setRecipients = (recipients: Recipient[]) => {
    setFormData((prev) => ({ ...prev, recipients }));
    if (errors.recipients) {
      setErrors((prev) => ({ ...prev, recipients: undefined }));
    }
  };

  const selectedValues = companyValues.filter((v) =>
    formData.valueIds.includes(v.id),
  );

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <Text as="h1" size={900} weight="bold">
          Send Recognition Card
        </Text>
        <Text
          as="p"
          size={300}
          style={{ color: tokens.colorNeutralForeground3 }}
        >
          Appreciate your colleagues and celebrate our company values
        </Text>
      </div>

      {successMessage && (
        <MessageBar intent="success" className={styles.message}>
          <MessageBarBody>
            <MessageBarTitle>Success</MessageBarTitle>
            {successMessage}
          </MessageBarBody>
        </MessageBar>
      )}

      {errorMessage && (
        <MessageBar intent="error" className={styles.message}>
          <MessageBarBody>
            <MessageBarTitle>Error</MessageBarTitle>
            {errorMessage}
          </MessageBarBody>
        </MessageBar>
      )}

      <form onSubmit={handleSubmit} className={styles.form}>
        <RecipientSelector
          selectedRecipients={formData.recipients}
          onChange={setRecipients}
          error={errors.recipients}
        />

        <Field
          label="Recognition Reason"
          required
          validationMessage={errors.recognitionReason}
          validationState={errors.recognitionReason ? "error" : "none"}
          className={styles.textareaField}
          hint={`${formData.recognitionReason.length}/500 characters`}
        >
          <Textarea
            value={formData.recognitionReason}
            onChange={(e) =>
              setFormData({ ...formData, recognitionReason: e.target.value })
            }
            placeholder="Describe why you're recognizing this person..."
            resize="vertical"
            rows={5}
          />
        </Field>

        <ValueSelector
          values={companyValues}
          selectedValueIds={formData.valueIds}
          onChange={(valueIds) => setFormData({ ...formData, valueIds })}
          error={errors.valueIds}
          loading={loadingValues}
        />

        <CardPreview
          recipients={formData.recipients}
          recognitionReason={formData.recognitionReason}
          selectedValues={selectedValues}
        />

        <div className={styles.actions}>
          <Button
            type="submit"
            appearance="primary"
            icon={<SendRegular />}
            disabled={loading}
          >
            {loading ? "Sending..." : "Send Recognition"}
          </Button>
          <Button
            type="button"
            appearance="secondary"
            onClick={handleReset}
            disabled={loading}
          >
            Reset
          </Button>
        </div>
      </form>
    </div>
  );
};
