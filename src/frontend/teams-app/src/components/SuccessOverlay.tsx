import React, { useEffect } from 'react';
import { makeStyles, tokens, Text } from '@fluentui/react-components';
import { CheckmarkCircle24Filled } from '@fluentui/react-icons';
import { useAnimation } from '../hooks/useAnimation';

const useStyles = makeStyles({
  overlay: {
    position: 'fixed',
    top: 0,
    left: 0,
    right: 0,
    bottom: 0,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    backgroundColor: tokens.colorNeutralBackgroundAlpha,
    zIndex: 1000,
    pointerEvents: 'none',
  },
  messageContainer: {
    backgroundColor: tokens.colorNeutralBackground1,
    borderRadius: tokens.borderRadiusXLarge,
    padding: tokens.spacingVerticalXXL,
    boxShadow: tokens.shadow64,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: tokens.spacingVerticalL,
    maxWidth: '400px',
    textAlign: 'center',
  },
  icon: {
    color: tokens.colorPaletteGreenForeground1,
    fontSize: '64px',
    lineHeight: 1,
  },
  message: {
    color: tokens.colorNeutralForeground1,
  },
  hidden: {
    display: 'none',
  },
});

export interface SuccessOverlayProps {
  /**
   * Whether the overlay is visible
   */
  visible: boolean;
  
  /**
   * Success message to display
   */
  message: string;
  
  /**
   * Callback when overlay auto-dismisses
   */
  onDismiss?: () => void;
  
  /**
   * Auto-dismiss duration in milliseconds (default: 5000)
   */
  dismissDelay?: number;
}

/**
 * SuccessOverlay Component
 * 
 * Displays a success message overlay with scale-in animation
 * and auto-dismisses after a specified delay.
 * 
 * Requirements: 4.2, 4.6
 */
export const SuccessOverlay: React.FC<SuccessOverlayProps> = ({
  visible,
  message,
  onDismiss,
  dismissDelay = 5000,
}) => {
  const styles = useStyles();
  const scaleInAnimation = useAnimation('scaleIn');
  const fadeOutAnimation = useAnimation('fadeOut');

  useEffect(() => {
    if (visible) {
      // Trigger scale-in animation when visible
      scaleInAnimation.trigger();

      // Set up auto-dismiss timer
      const dismissTimer = setTimeout(() => {
        // Trigger fade-out animation before dismissing
        fadeOutAnimation.trigger();
        
        // Call onDismiss after fade-out completes
        setTimeout(() => {
          onDismiss?.();
        }, 300); // Match fadeOut animation duration
      }, dismissDelay);

      return () => {
        clearTimeout(dismissTimer);
      };
    }
  }, [visible, dismissDelay, onDismiss, scaleInAnimation, fadeOutAnimation]);

  if (!visible) {
    return null;
  }

  return (
    <div 
      className={`${styles.overlay} ${fadeOutAnimation.animationClass}`}
      role="alert"
      aria-live="polite"
    >
      <div className={`${styles.messageContainer} ${scaleInAnimation.animationClass}`}>
        <CheckmarkCircle24Filled className={styles.icon} />
        <Text size={600} weight="semibold" className={styles.message}>
          {message}
        </Text>
      </div>
    </div>
  );
};
