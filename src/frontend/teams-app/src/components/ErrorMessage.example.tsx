/**
 * ErrorMessage Component Examples
 * 
 * This file demonstrates various use cases of the ErrorMessage component
 */

import React, { useState } from "react";
import { Button, makeStyles, tokens } from "@fluentui/react-components";
import { ErrorMessage } from "./ErrorMessage";

const useStyles = makeStyles({
  container: {
    padding: tokens.spacingVerticalXL,
    maxWidth: "800px",
    margin: "0 auto",
  },
  section: {
    marginBottom: tokens.spacingVerticalXXL,
  },
  title: {
    fontSize: tokens.fontSizeBase500,
    fontWeight: tokens.fontWeightSemibold,
    marginBottom: tokens.spacingVerticalM,
  },
  controls: {
    display: "flex",
    gap: tokens.spacingHorizontalM,
    marginBottom: tokens.spacingVerticalL,
  },
  description: {
    color: tokens.colorNeutralForeground3,
    marginBottom: tokens.spacingVerticalM,
  },
});

export const ErrorMessageExamples: React.FC = () => {
  const styles = useStyles();
  
  // Example 1: Simple error message with auto-dismiss
  const [showSimpleError, setShowSimpleError] = useState(false);
  
  // Example 2: Retryable error with retry callback
  const [showRetryableError, setShowRetryableError] = useState(false);
  const [retryCount, setRetryCount] = useState(0);
  
  // Example 3: Persistent error (no auto-dismiss)
  const [showPersistentError, setShowPersistentError] = useState(false);
  
  // Example 4: Multiple errors
  const [errors, setErrors] = useState<string[]>([]);

  const handleRetry = () => {
    setRetryCount((prev) => prev + 1);
    console.log(`Retry attempt ${retryCount + 1}`);
    // Simulate retry logic
    setTimeout(() => {
      setShowRetryableError(false);
      alert("Retry successful!");
    }, 1000);
  };

  const addMultipleErrors = () => {
    setErrors([
      "Failed to load user data",
      "Network connection timeout",
      "Invalid authentication token",
    ]);
  };

  const removeError = (index: number) => {
    setErrors((prev) => prev.filter((_, i) => i !== index));
  };

  return (
    <div className={styles.container}>
      <h1>ErrorMessage Component Examples</h1>

      {/* Example 1: Simple Error */}
      <div className={styles.section}>
        <h2 className={styles.title}>1. Simple Error (Auto-dismiss)</h2>
        <p className={styles.description}>
          A basic error message that automatically dismisses after 10 seconds.
        </p>
        <div className={styles.controls}>
          <Button onClick={() => setShowSimpleError(true)}>
            Show Simple Error
          </Button>
        </div>
        {showSimpleError && (
          <ErrorMessage
            message="Something went wrong. Please try again later."
            onDismiss={() => setShowSimpleError(false)}
            visible={showSimpleError}
          />
        )}
      </div>

      {/* Example 2: Retryable Error */}
      <div className={styles.section}>
        <h2 className={styles.title}>2. Retryable Error</h2>
        <p className={styles.description}>
          An error with a retry button. Does not auto-dismiss by default.
        </p>
        <div className={styles.controls}>
          <Button onClick={() => setShowRetryableError(true)}>
            Show Retryable Error
          </Button>
        </div>
        {showRetryableError && (
          <ErrorMessage
            title="Network Error"
            message="Failed to connect to the server. Please check your connection and try again."
            retryable
            onRetry={handleRetry}
            onDismiss={() => setShowRetryableError(false)}
            visible={showRetryableError}
          />
        )}
      </div>

      {/* Example 3: Persistent Error */}
      <div className={styles.section}>
        <h2 className={styles.title}>3. Persistent Error (No Auto-dismiss)</h2>
        <p className={styles.description}>
          An error that requires user action to dismiss.
        </p>
        <div className={styles.controls}>
          <Button onClick={() => setShowPersistentError(true)}>
            Show Persistent Error
          </Button>
        </div>
        {showPersistentError && (
          <ErrorMessage
            title="Critical Error"
            message="Unable to save your changes. Please contact support if this issue persists."
            autoDismiss={false}
            onDismiss={() => setShowPersistentError(false)}
            visible={showPersistentError}
          />
        )}
      </div>

      {/* Example 4: Multiple Errors */}
      <div className={styles.section}>
        <h2 className={styles.title}>4. Multiple Errors (Stacked)</h2>
        <p className={styles.description}>
          Multiple error messages displayed with proper spacing.
        </p>
        <div className={styles.controls}>
          <Button onClick={addMultipleErrors}>Show Multiple Errors</Button>
          <Button onClick={() => setErrors([])}>Clear All</Button>
        </div>
        {errors.map((error, index) => (
          <ErrorMessage
            key={index}
            message={error}
            onDismiss={() => removeError(index)}
            visible={true}
            dismissDelay={15000}
          />
        ))}
      </div>

      {/* Example 5: Custom Dismiss Delay */}
      <div className={styles.section}>
        <h2 className={styles.title}>5. Custom Dismiss Delay (3 seconds)</h2>
        <p className={styles.description}>
          An error that auto-dismisses after 3 seconds instead of the default 10.
        </p>
        <div className={styles.controls}>
          <Button onClick={() => setShowSimpleError(true)}>
            Show Quick Dismiss Error
          </Button>
        </div>
        {showSimpleError && (
          <ErrorMessage
            message="This error will disappear in 3 seconds."
            dismissDelay={3000}
            onDismiss={() => setShowSimpleError(false)}
            visible={showSimpleError}
          />
        )}
      </div>
    </div>
  );
};

export default ErrorMessageExamples;
