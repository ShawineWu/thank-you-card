import { useState, useCallback, useEffect, useRef } from 'react';

/**
 * Animation configuration interface
 */
export interface AnimationConfig {
  duration: number;        // Animation duration in milliseconds
  easing: string;          // CSS easing function
  delay?: number;          // Delay before animation starts
  fillMode?: 'forwards' | 'backwards' | 'both' | 'none';
}

/**
 * Predefined animation presets for common animations
 */
export const AnimationPresets = {
  fadeIn: {
    duration: 300,
    easing: 'ease-in-out',
    fillMode: 'forwards' as const,
  },
  fadeOut: {
    duration: 300,
    easing: 'ease-in-out',
    fillMode: 'forwards' as const,
  },
  slideIn: {
    duration: 400,
    easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
    fillMode: 'forwards' as const,
  },
  slideOut: {
    duration: 400,
    easing: 'cubic-bezier(0.4, 0, 0.2, 1)',
    fillMode: 'forwards' as const,
  },
  scaleIn: {
    duration: 300,
    easing: 'cubic-bezier(0.34, 1.56, 0.64, 1)',
    fillMode: 'forwards' as const,
  },
  scaleOut: {
    duration: 300,
    easing: 'ease-in',
    fillMode: 'forwards' as const,
  },
  shake: {
    duration: 500,
    easing: 'ease-in-out',
    fillMode: 'forwards' as const,
  },
  pulse: {
    duration: 1000,
    easing: 'ease-in-out',
    fillMode: 'forwards' as const,
  },
  flip: {
    duration: 600,
    easing: 'ease-in-out',
    fillMode: 'forwards' as const,
  },
  bounce: {
    duration: 600,
    easing: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)',
    fillMode: 'forwards' as const,
  },
  stagger: {
    duration: 300,
    easing: 'ease-out',
    fillMode: 'forwards' as const,
  },
} as const;

export type AnimationType = keyof typeof AnimationPresets;

/**
 * Return type for useAnimation hook
 */
export interface UseAnimationReturn {
  trigger: () => void;
  reset: () => void;
  isAnimating: boolean;
  animationClass: string;
}

/**
 * Custom hook for managing animations with GPU acceleration
 * 
 * @param animationType - Type of animation from AnimationPresets
 * @param options - Optional configuration to override preset defaults
 * @returns Animation control interface
 */
export function useAnimation(
  animationType: AnimationType,
  options?: Partial<AnimationConfig>
): UseAnimationReturn {
  const [isAnimating, setIsAnimating] = useState(false);
  const timeoutRef = useRef<number | null>(null);
  
  const config = {
    ...AnimationPresets[animationType],
    ...options,
  };

  const trigger = useCallback(() => {
    // Clear any existing timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    setIsAnimating(true);

    // Calculate total animation time including delay
    const totalDuration = config.duration + (config.delay || 0);

    // Reset animation state after completion
    timeoutRef.current = window.setTimeout(() => {
      setIsAnimating(false);
      timeoutRef.current = null;
    }, totalDuration);
  }, [config.duration, config.delay]);

  const reset = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
    setIsAnimating(false);
  }, []);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  // Generate animation class name
  const animationClass = isAnimating ? `animate-${animationType}` : '';

  return {
    trigger,
    reset,
    isAnimating,
    animationClass,
  };
}

/**
 * Hook for detecting reduced motion preference
 */
export function useReducedMotion(): boolean {
  const [prefersReducedMotion, setPrefersReducedMotion] = useState(false);

  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    setPrefersReducedMotion(mediaQuery.matches);

    const handleChange = (event: MediaQueryListEvent) => {
      setPrefersReducedMotion(event.matches);
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  return prefersReducedMotion;
}
