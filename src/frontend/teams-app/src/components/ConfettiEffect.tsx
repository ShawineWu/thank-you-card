import { useEffect, useRef, useCallback } from 'react';
import confetti from 'canvas-confetti';
import { useThemeColors } from '../hooks/useThemeColors';
import { SuccessAnimationConfig } from '../config/animations';

/**
 * Confetti configuration interface
 */
export interface ConfettiConfig {
  /** Number of confetti particles to generate */
  particleCount: number;
  /** Spread angle in degrees (0-360) */
  spread: number;
  /** Origin point for confetti (x: 0-1, y: 0-1) */
  origin: { x: number; y: number };
  /** Array of colors for confetti particles */
  colors: string[];
  /** Duration of the confetti effect in milliseconds */
  duration: number;
}

/**
 * Props for ConfettiEffect component
 */
export interface ConfettiEffectProps {
  /** Whether the confetti effect is active */
  active: boolean;
  /** Callback when confetti animation completes */
  onComplete?: () => void;
  /** Optional custom configuration (overrides defaults) */
  config?: Partial<ConfettiConfig>;
}

/**
 * Default confetti configuration
 */
const DEFAULT_CONFIG: ConfettiConfig = {
  particleCount: SuccessAnimationConfig.CONFETTI_PARTICLE_COUNT,
  spread: 70,
  origin: { x: 0.5, y: 0.5 },
  colors: [],
  duration: SuccessAnimationConfig.CONFETTI_DURATION,
};

/**
 * ConfettiEffect Component
 * 
 * Displays a celebratory confetti animation using canvas-confetti library.
 * The component is theme-aware and uses colors from the current Fluent UI theme.
 * 
 * Features:
 * - Configurable particle count, spread, and duration
 * - Theme-aware color selection
 * - Multiple confetti bursts for enhanced effect
 * - Automatic cleanup and fade-out
 * 
 * Requirements: 1.1, 4.1, 4.2, 10.4
 * 
 * @example
 * ```tsx
 * <ConfettiEffect 
 *   active={showConfetti}
 *   onComplete={() => setShowConfetti(false)}
 *   config={{ particleCount: 200, spread: 90 }}
 * />
 * ```
 */
export const ConfettiEffect: React.FC<ConfettiEffectProps> = ({
  active,
  onComplete,
  config: customConfig,
}) => {
  const themeColors = useThemeColors();
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const animationFrameRef = useRef<number | null>(null);

  /**
   * Get theme-aware confetti colors
   * Extracts actual color values from CSS variables
   */
  const getThemeColors = useCallback((): string[] => {
    // If custom colors provided, use them
    if (customConfig?.colors && customConfig.colors.length > 0) {
      return customConfig.colors;
    }

    // Extract colors from theme tokens
    // These are CSS variables, so we need to get computed values
    const colors: string[] = [];
    
    // Helper to extract color from CSS variable
    const extractColor = (cssVar: string): string | null => {
      if (typeof window === 'undefined') return null;
      
      // Create a temporary element to compute the color
      const temp = document.createElement('div');
      temp.style.color = cssVar;
      document.body.appendChild(temp);
      const computed = window.getComputedStyle(temp).color;
      document.body.removeChild(temp);
      
      return computed;
    };

    // Use theme colors for confetti
    const primaryColor = extractColor(themeColors.primary);
    const secondaryColor = extractColor(themeColors.secondary);
    const successColor = extractColor(themeColors.success);
    const warningColor = extractColor(themeColors.warning);
    
    if (primaryColor) colors.push(primaryColor);
    if (secondaryColor) colors.push(secondaryColor);
    if (successColor) colors.push(successColor);
    if (warningColor) colors.push(warningColor);
    
    // Fallback colors only if extraction fails (should not happen in normal operation)
    if (colors.length === 0) {
      console.warn('Failed to extract theme colors for confetti, using fallback colors');
      colors.push('#0078D4', '#50E6FF', '#00CC6A', '#FFB900');
    }
    
    return colors;
  }, [themeColors, customConfig?.colors]);

  /**
   * Fire confetti burst
   */
  const fireConfetti = useCallback(() => {
    const mergedConfig: ConfettiConfig = {
      ...DEFAULT_CONFIG,
      ...customConfig,
      colors: getThemeColors(),
    };

    // Fire multiple bursts for better effect
    const burstCount = 3;
    const burstDelay = 150;

    for (let i = 0; i < burstCount; i++) {
      setTimeout(() => {
        confetti({
          particleCount: Math.floor(mergedConfig.particleCount / burstCount),
          spread: mergedConfig.spread,
          origin: mergedConfig.origin,
          colors: mergedConfig.colors,
          startVelocity: 30,
          gravity: 0.8,
          drift: 0,
          ticks: 200,
          scalar: 1,
          zIndex: 1001, // Above overlay
        });
      }, i * burstDelay);
    }

    // Schedule completion callback
    timeoutRef.current = setTimeout(() => {
      onComplete?.();
    }, mergedConfig.duration);
  }, [customConfig, getThemeColors, onComplete]);

  /**
   * Trigger confetti when active becomes true
   */
  useEffect(() => {
    if (active) {
      fireConfetti();
    }

    // Cleanup on unmount or when active changes
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
        timeoutRef.current = null;
      }
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
        animationFrameRef.current = null;
      }
    };
  }, [active, fireConfetti]);

  // This component doesn't render anything visible
  // The confetti is rendered on a canvas overlay by the library
  return null;
};
