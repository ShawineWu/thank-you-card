import { useMemo } from 'react';
import { tokens } from '@fluentui/react-components';

/**
 * Theme colors interface providing type-safe access to theme tokens
 */
export interface ThemeColors {
  primary: string;
  secondary: string;
  success: string;
  warning: string;
  error: string;
  background: string;
  surface: string;
  text: string;
  textSecondary: string;
  border: string;
  shadow: string;
}

/**
 * Custom hook for accessing theme colors in a type-safe manner
 * 
 * Provides access to Fluent UI theme tokens (CSS variables) that automatically
 * adapt to the current theme (light, dark, or high contrast) set by FluentProvider.
 * 
 * The returned colors are CSS variable references that will update automatically
 * when the theme changes, ensuring consistent styling across all theme modes.
 * 
 * @returns ThemeColors object with semantic color tokens as CSS variables
 * 
 * @example
 * ```tsx
 * function MyComponent() {
 *   const colors = useThemeColors();
 *   
 *   return (
 *     <div style={{ 
 *       backgroundColor: colors.surface,
 *       color: colors.text,
 *       borderColor: colors.border
 *     }}>
 *       Content
 *     </div>
 *   );
 * }
 * ```
 */
export function useThemeColors(): ThemeColors {
  return useMemo(() => ({
    // Primary brand color
    primary: tokens.colorBrandBackground,
    
    // Secondary/accent color
    secondary: tokens.colorBrandBackgroundHover,
    
    // Success state color (green)
    success: tokens.colorPaletteGreenBackground3,
    
    // Warning state color (yellow/orange)
    warning: tokens.colorPaletteYellowBackground3,
    
    // Error state color (red)
    error: tokens.colorPaletteRedBackground3,
    
    // Main background color
    background: tokens.colorNeutralBackground1,
    
    // Surface/card background color
    surface: tokens.colorNeutralBackground2,
    
    // Primary text color
    text: tokens.colorNeutralForeground1,
    
    // Secondary/muted text color
    textSecondary: tokens.colorNeutralForeground2,
    
    // Border color
    border: tokens.colorNeutralStroke1,
    
    // Shadow color for elevation
    shadow: tokens.colorNeutralShadowAmbient,
  }), []);
}
