import React, { useState, useEffect, useRef } from "react";
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
  Spinner,
} from "@fluentui/react-components";
import { SendRegular } from "@fluentui/react-icons";
import * as microsoftTeams from "@microsoft/teams-js";
import { RecipientSelector } from "./RecipientSelector";
import { ValueSelector } from "./ValueSelector";
import { CardPreview } from "./CardPreview";
import { CharacterCounter } from "./CharacterCounter";
import { SuccessOverlay } from "./SuccessOverlay";
import { ConfettiEffect } from "./ConfettiEffect";
import { ErrorMessage } from "./ErrorMessage";
import { useAnimation } from "../hooks/useAnimation";
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
  fieldWithError: {
    position: "relative",
  },
  errorMessage: {
    color: tokens.colorPaletteRedForeground1,
    fontSize: tokens.fontSizeBase200,
    marginTop: tokens.spacingVerticalXXS,
    display: "block",
  },
  errorMessageEnter: {
    animation: "fadeIn 300ms ease-in-out forwards",
  },
  errorMessageExit: {
    animation: "fadeOut 300ms ease-in-out forwards",
  },
  submitButton: {
    position: "relative",
  },
  loadingSpinner: {
    display: "inline-flex",
    alignItems: "center",
    gap: tokens.spacingHorizontalS,
  },
  fieldResetting: {
    animation: "fadeOut 300ms ease-out forwards",
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
  const [apiError, setApiError] = useState<{
    message: string;
    retryable: boolean;
    operation?: () => Promise<void>;
  } | null>(null);
  const [companyValues, setCompanyValues] = useState<CompanyValue[]>([]);
  const [loadingValues, setLoadingValues] = useState(true);
  const [isInTeamsTask, setIsInTeamsTask] = useState(false);
  const [showSuccessOverlay, setShowSuccessOverlay] = useState(false);
  const [showConfetti, setShowConfetti] = useState(false);
  const [isResetting, setIsResetting] = useState(false);
  
  // Animation hooks for validation errors
  const recipientShake = useAnimation("shake");
  const reasonShake = useAnimation("shake");
  const valuesShake = useAnimation("shake");
  const submitButtonPulse = useAnimation("pulse");
  
  // Refs for field elements to apply shake animation
  const recipientFieldRef = useRef<HTMLDivElement>(null);
  const reasonFieldRef = useRef<HTMLDivElement>(null);
  const valuesFieldRef = useRef<HTMLDivElement>(null);
  
  // Check if form is valid
  const isFormValid = 
    formData.recipients.length > 0 &&
    formData.recognitionReason.trim().length >= 10 &&
    formData.recognitionReason.length <= 500 &&
    formData.valueIds.length > 0;

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
        setApiError(null);
        const values = await api.getCompanyValues();
        setCompanyValues(values);
      } catch (error: any) {
        const errorMessage = error.message || "Failed to load company values";
        setApiError({
          message: errorMessage,
          retryable: true,
          operation: fetchValues,
        });
      } finally {
        setLoadingValues(false);
      }
    };
    fetchValues();
  }, []);
  
  // Trigger pulse animation when form becomes valid
  useEffect(() => {
    if (isFormValid && !loading) {
      submitButtonPulse.trigger();
    }
  }, [isFormValid, loading, submitButtonPulse]);

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
    
    // Trigger shake animation for all invalid fields
    if (newErrors.recipients) {
      recipientShake.trigger();
    }
    if (newErrors.recognitionReason) {
      reasonShake.trigger();
    }
    if (newErrors.valueIds) {
      valuesShake.trigger();
    }
    
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSuccessMessage("");
    setApiError(null);

    if (!validate()) {
      return;
    }

    const submitCard = async () => {
      try {
        setLoading(true);
        setApiError(null);
        const cardRequest = {
          recipients: formData.recipients,
          recognitionReason: formData.recognitionReason.trim(),
          valueIds: formData.valueIds,
        };

        await api.createCard(cardRequest);

        setSuccessMessage("Recognition card sent successfully! 🎉");
        
        // Trigger confetti and success overlay
        setShowConfetti(true);
        setShowSuccessOverlay(true);

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

        // Reset form if not in a task (will be handled by form reset animation)
        // Delay reset to allow success animation to play
        setTimeout(() => {
          resetFormWithAnimation();
        }, 5000); // Wait for success overlay to auto-dismiss
      } catch (error: any) {
        const errorMessage = error.message || "Failed to send recognition card";
        setApiError({
          message: errorMessage,
          retryable: true,
          operation: submitCard,
        });
      } finally {
        if (!isInTeamsTask) {
          setLoading(false);
        }
      }
    };

    await submitCard();
  };

  const handleReset = () => {
    resetFormWithAnimation();
  };

  const resetFormWithAnimation = () => {
    setIsResetting(true);
    
    // Staggered field clearing animation
    // Clear recipients first
    setTimeout(() => {
      setFormData((prev) => ({ ...prev, recipients: [] }));
    }, 100);
    
    // Clear recognition reason second
    setTimeout(() => {
      setFormData((prev) => ({ ...prev, recognitionReason: "" }));
    }, 200);
    
    // Clear values third
    setTimeout(() => {
      setFormData((prev) => ({ ...prev, valueIds: [] }));
    }, 300);
    
    // Clear errors and messages
    setTimeout(() => {
      setErrors({});
      setSuccessMessage("");
      setApiError(null);
      setShowSuccessOverlay(false);
      setShowConfetti(false);
      setIsResetting(false);
    }, 400);
  };

  const handleRetry = async () => {
    if (apiError?.operation) {
      await apiError.operation();
    }
  };

  const handleErrorDismiss = () => {
    setApiError(null);
  };

  const handleSuccessOverlayDismiss = () => {
    setShowSuccessOverlay(false);
    setSuccessMessage("");
  };

  const handleConfettiComplete = () => {
    setShowConfetti(false);
  };

  const setRecipients = (recipients: Recipient[]) => {
    setFormData((prev) => ({ ...prev, recipients }));
    if (errors.recipients && recipients.length > 0) {
      setErrors((prev) => ({ ...prev, recipients: undefined }));
    }
  };
  
  const handleReasonChange = (value: string) => {
    setFormData((prev) => ({ ...prev, recognitionReason: value }));
    
    // Clear error if field becomes valid
    if (errors.recognitionReason) {
      if (value.trim().length >= 10 && value.length <= 500) {
        setErrors((prev) => ({ ...prev, recognitionReason: undefined }));
      }
    }
  };
  
  const handleValuesChange = (valueIds: string[]) => {
    setFormData((prev) => ({ ...prev, valueIds }));
    if (errors.valueIds && valueIds.length > 0) {
      setErrors((prev) => ({ ...prev, valueIds: undefined }));
    }
  };

  const selectedValues = companyValues.filter((v) =>
    formData.valueIds.includes(v.id),
  );

  return (
    <div className={styles.container}>
      {/* Confetti Effect */}
      <ConfettiEffect 
        active={showConfetti} 
        onComplete={handleConfettiComplete}
      />
      
      {/* Success Overlay */}
      <SuccessOverlay
        visible={showSuccessOverlay}
        message="Recognition card sent successfully! 🎉"
        onDismiss={handleSuccessOverlayDismiss}
        dismissDelay={5000}
      />
      
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

      {apiError && (
        <ErrorMessage
          title={apiError.retryable ? "Network Error" : "Error"}
          message={apiError.message}
          retryable={apiError.retryable}
          onRetry={handleRetry}
          onDismiss={handleErrorDismiss}
          visible={true}
        />
      )}

      <form onSubmit={handleSubmit} className={styles.form}>
        <div ref={recipientFieldRef} className={`${recipientShake.animationClass} ${isResetting ? styles.fieldResetting : ''}`}>
          <RecipientSelector
            selectedRecipients={formData.recipients}
            onChange={setRecipients}
            error={errors.recipients}
          />
        </div>

        <div ref={reasonFieldRef} className={`${reasonShake.animationClass} ${isResetting ? styles.fieldResetting : ''}`}>
          <Field
            label="Recognition Reason"
            required
            validationMessage={errors.recognitionReason}
            validationState={errors.recognitionReason ? "error" : "none"}
            className={styles.textareaField}
            hint={
              <CharacterCounter
                current={formData.recognitionReason.length}
                max={500}
                warningThreshold={0.8}
                errorThreshold={1.0}
              />
            }
          >
            <Textarea
              value={formData.recognitionReason}
              onChange={(e) => handleReasonChange(e.target.value)}
              placeholder="Describe why you're recognizing this person..."
              resize="vertical"
              rows={5}
            />
          </Field>
        </div>

        <div ref={valuesFieldRef} className={`${valuesShake.animationClass} ${isResetting ? styles.fieldResetting : ''}`}>
          <ValueSelector
            values={companyValues}
            selectedValueIds={formData.valueIds}
            onChange={handleValuesChange}
            error={errors.valueIds}
            loading={loadingValues}
          />
        </div>

        <CardPreview
          recipients={formData.recipients}
          recognitionReason={formData.recognitionReason}
          selectedValues={selectedValues}
        />

        <div className={styles.actions}>
          <Button
            type="submit"
            appearance="primary"
            icon={loading ? undefined : <SendRegular />}
            disabled={loading || !isFormValid}
            className={`${styles.submitButton} ${isFormValid && !loading ? submitButtonPulse.animationClass : ''}`}
          >
            {loading ? (
              <span className={styles.loadingSpinner}>
                <Spinner size="tiny" />
                <span>Sending...</span>
              </span>
            ) : (
              "Send Recognition"
            )}
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
