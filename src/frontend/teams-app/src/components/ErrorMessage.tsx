import React, { useEffect, useState, useCallback } from "react";
import {
  MessageBar,
  MessageBarBody,
  MessageBarTitle,
  Button,
  makeStyles,
  tokens,
} from "@fluentui/react-components";
import { DismissRegular, ArrowClockwiseRegular } from "@fluentui/react-icons";
import { ErrorAnimationConfig } from "../config/animations";

const useStyles = makeStyles({
  errorContainer: {
    position: "relative",
    marginBottom: tokens.spacingVerticalM,
    opacity: 0,
    transform: "translateY(-10px)",
    animation: "slideIn 300ms ease-out forwards",
    willChange: "transform, opacity",
  },
  errorContainerExit: {
    animation: "fadeOut 300ms ease-out forwards",
  },
  messageBarContent: {
    display: "flex",
    alignItems: "center",
    justifyContent: "space-between",
    width: "100%",
  },
  messageText: {
    flex: 1,
  },
  actions: {
    display: "flex",
    gap: tokens.spacingHorizontalS,
    marginLeft: tokens.spacingHorizontalM,
  },
  retryButton: {
    minWidth: "auto",
  },
  dismissButton: {
    minWidth: "auto",
  },
});

export interface ErrorMessageProps {
  /**
   * Error message to display
   */
  message: string;
  
  /**
   * Optional error title (defaults to "Error")
   */
  title?: string;
  
  /**
   * Whether the error is retryable
   */
  retryable?: boolean;
  
  /**
   * Callback function when retry button is clicked
   */
  onRetry?: () => void;
  
  /**
   * Whether to auto-dismiss the error after a delay
   * Defaults to true for non-retryable errors, false for retryable errors
   */
  autoDismiss?: boolean;
  
  /**
   * Auto-dismiss delay in milliseconds (defaults to 10000ms)
   */
  dismissDelay?: number;
  
  /**
   * Callback function when error is dismissed
   */
  onDismiss?: () => void;
  
  /**
   * Whether the error message is visible
   */
  visible?: boolean;
}

export const ErrorMessage: React.FC<ErrorMessageProps> = ({
  message,
  title = "Error",
  retryable = false,
  onRetry,
  autoDismiss,
  dismissDelay = ErrorAnimationConfig.ERROR_MESSAGE_DURATION,
  onDismiss,
  visible = true,
}) => {
  const styles = useStyles();
  const [isExiting, setIsExiting] = useState(false);
  const [prevVisible, setPrevVisible] = useState(visible);
  
  // Determine auto-dismiss behavior
  const shouldAutoDismiss = autoDismiss !== undefined 
    ? autoDismiss 
    : !retryable; // Auto-dismiss non-retryable errors by default

  // Reset exiting state when visibility changes from false to true
  if (visible !== prevVisible) {
    setPrevVisible(visible);
    if (visible && isExiting) {
      setIsExiting(false);
    }
  }

  const handleDismiss = useCallback(() => {
    setIsExiting(true);
    setTimeout(() => {
      onDismiss?.();
    }, ErrorAnimationConfig.ERROR_FADE_OUT_DURATION);
  }, [onDismiss]);

  useEffect(() => {
    if (!visible || !shouldAutoDismiss) {
      return;
    }

    const timer = setTimeout(() => {
      handleDismiss();
    }, dismissDelay);

    return () => clearTimeout(timer);
  }, [visible, shouldAutoDismiss, dismissDelay, handleDismiss]);

  const handleRetry = () => {
    if (onRetry) {
      onRetry();
    }
  };

  if (!visible || isExiting) {
    return null;
  }

  return (
    <div
      className={`${styles.errorContainer} ${isExiting ? styles.errorContainerExit : ""}`}
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      <MessageBar intent="error">
        <MessageBarBody>
          <div className={styles.messageBarContent}>
            <div className={styles.messageText}>
              <MessageBarTitle>{title}</MessageBarTitle>
              {message}
            </div>
            <div className={styles.actions}>
              {retryable && onRetry && (
                <Button
                  appearance="subtle"
                  icon={<ArrowClockwiseRegular />}
                  onClick={handleRetry}
                  className={styles.retryButton}
                  aria-label="Retry"
                >
                  Retry
                </Button>
              )}
              <Button
                appearance="subtle"
                icon={<DismissRegular />}
                onClick={handleDismiss}
                className={styles.dismissButton}
                aria-label="Dismiss error"
              />
            </div>
          </div>
        </MessageBarBody>
      </MessageBar>
    </div>
  );
};
