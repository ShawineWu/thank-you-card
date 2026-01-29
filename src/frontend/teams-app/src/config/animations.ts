/**
 * Animation Configuration
 * 
 * Centralized configuration for all animations in the application
 */

/**
 * Animation duration constants (in milliseconds)
 */
export const AnimationDuration = {
  INSTANT: 0,
  FAST: 150,
  NORMAL: 300,
  SLOW: 500,
  VERY_SLOW: 1000,
} as const;

/**
 * Animation easing functions
 */
export const AnimationEasing = {
  LINEAR: 'linear',
  EASE: 'ease',
  EASE_IN: 'ease-in',
  EASE_OUT: 'ease-out',
  EASE_IN_OUT: 'ease-in-out',
  EASE_IN_CUBIC: 'cubic-bezier(0.4, 0, 1, 1)',
  EASE_OUT_CUBIC: 'cubic-bezier(0, 0, 0.2, 1)',
  EASE_IN_OUT_CUBIC: 'cubic-bezier(0.4, 0, 0.2, 1)',
  EASE_IN_BACK: 'cubic-bezier(0.6, -0.28, 0.735, 0.045)',
  EASE_OUT_BACK: 'cubic-bezier(0.175, 0.885, 0.32, 1.275)',
  EASE_IN_OUT_BACK: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)',
  BOUNCE: 'cubic-bezier(0.68, -0.55, 0.265, 1.55)',
} as const;

/**
 * Stagger animation delays (in milliseconds)
 */
export const StaggerDelay = {
  ITEM_1: 50,
  ITEM_2: 100,
  ITEM_3: 150,
  ITEM_4: 200,
  ITEM_5: 250,
  ITEM_6: 300,
} as const;

/**
 * Success animation configuration
 */
export const SuccessAnimationConfig = {
  CONFETTI_DURATION: 3000,
  CONFETTI_PARTICLE_COUNT: 150,
  SUCCESS_MESSAGE_DURATION: 5000,
  SUCCESS_MESSAGE_FADE_OUT: 500,
} as const;

/**
 * Loading animation configuration
 */
export const LoadingAnimationConfig = {
  SKELETON_SHIMMER_DURATION: 2000,
  SPINNER_ROTATION_DURATION: 1000,
  FADE_IN_DELAY: 200,
} as const;

/**
 * Form animation configuration
 */
export const FormAnimationConfig = {
  VALIDATION_ERROR_SHAKE_DURATION: 500,
  FIELD_FOCUS_TRANSITION: 200,
  SUBMIT_BUTTON_PULSE_DURATION: 1000,
  CHARACTER_COUNTER_TRANSITION: 200,
} as const;

/**
 * Preview animation configuration
 */
export const PreviewAnimationConfig = {
  UPDATE_DEBOUNCE_DELAY: 300,
  CONTENT_TRANSITION_DURATION: 300,
  HIGHLIGHT_DURATION: 1000,
  VIEW_FLIP_DURATION: 600,
} as const;

/**
 * History animation configuration
 */
export const HistoryAnimationConfig = {
  CARD_STAGGER_DELAY: 100,
  CARD_HOVER_DURATION: 200,
  CARD_EXPAND_DURATION: 400,
  NEW_CARD_SLIDE_DURATION: 500,
} as const;

/**
 * Error animation configuration
 */
export const ErrorAnimationConfig = {
  ERROR_MESSAGE_DURATION: 10000,
  ERROR_FADE_IN_DURATION: 300,
  ERROR_FADE_OUT_DURATION: 300,
  RETRY_BUTTON_PULSE_DURATION: 1000,
} as const;

/**
 * Accessibility configuration
 */
export const AccessibilityConfig = {
  FOCUS_RING_WIDTH: 3,
  FOCUS_RING_COLOR: 'rgba(0, 120, 212, 0.4)',
  MIN_TOUCH_TARGET_SIZE: 44, // pixels
  REDUCED_MOTION_DURATION: 0.01, // milliseconds
} as const;

/**
 * Helper function to get animation class name
 */
export function getAnimationClass(
  animationType: string,
  isActive: boolean
): string {
  return isActive ? `animate-${animationType}` : '';
}

/**
 * Helper function to get stagger delay class
 */
export function getStaggerDelayClass(index: number): string {
  const delayIndex = Math.min(index + 1, 6);
  return `stagger-delay-${delayIndex}`;
}

/**
 * Helper function to check if reduced motion is preferred
 */
export function shouldReduceMotion(): boolean {
  if (typeof window === 'undefined') return false;
  
  const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
  return mediaQuery.matches;
}

/**
 * Helper function to get animation duration based on reduced motion preference
 */
export function getAnimationDuration(
  normalDuration: number,
  respectReducedMotion: boolean = true
): number {
  if (respectReducedMotion && shouldReduceMotion()) {
    return AccessibilityConfig.REDUCED_MOTION_DURATION;
  }
  return normalDuration;
}
