# Theme Adaptation Manual Verification Guide

## Overview

This document provides instructions for manually verifying that all components properly adapt to different Teams theme modes (light, dark, and high contrast).

**Requirements Validated:** 10.1, 10.2

## Prerequisites

1. Teams app is running in development mode (`npm run dev`)
2. Access to Microsoft Teams (web or desktop client)
3. Ability to change Teams theme settings

## Theme Modes to Test

### 1. Light Theme (Default)
### 2. Dark Theme
### 3. High Contrast Theme

## How to Change Teams Theme

1. Open Microsoft Teams
2. Click on your profile picture (top right)
3. Select **Settings**
4. Go to **General** tab
5. Under **Theme**, select:
   - **Default** for Light theme
   - **Dark** for Dark theme
   - **High contrast** for High Contrast theme

## Components to Verify

### ✅ CardPreview Component

**Elements to check:**
- [ ] Card background uses theme token (`colorNeutralBackground1`)
- [ ] Content header background uses theme token (`colorNeutralBackground2`)
- [ ] Envelope container uses brand color (`colorBrandBackground`)
- [ ] Envelope text uses on-brand foreground (`colorNeutralForegroundOnBrand`)
- [ ] Open button uses theme colors for background and border
- [ ] All text colors adapt to theme

**Expected behavior:**
- Light theme: White/light gray backgrounds, dark text
- Dark theme: Dark backgrounds, light text
- High contrast: Maximum contrast between text and background

---

### ✅ ConfettiEffect Component

**Elements to check:**
- [ ] Confetti particles use colors extracted from current theme
- [ ] Primary color from theme is used
- [ ] Secondary color from theme is used
- [ ] Success color from theme is used
- [ ] Warning color from theme is used

**Expected behavior:**
- Confetti colors should match the current theme's color palette
- No hardcoded colors should be visible
- Colors should be vibrant and celebratory in all themes

**How to test:**
1. Fill out the card form completely
2. Submit the card
3. Observe the confetti colors
4. Change theme and repeat

---

### ✅ CharacterCounter Component

**Elements to check:**
- [ ] Normal state uses `colorNeutralForeground3`
- [ ] Warning state uses `colorPaletteYellowForeground1`
- [ ] Error state uses `colorPaletteRedForeground1`

**Expected behavior:**
- All states should be clearly visible in all themes
- Color transitions should be smooth
- Text should maintain readability

**How to test:**
1. Type in the recognition reason field
2. Observe counter color at different character counts:
   - 0-400 chars: Normal (gray)
   - 400-500 chars: Warning (yellow/orange)
   - 500+ chars: Error (red)

---

### ✅ SkeletonLoader Component

**Elements to check:**
- [ ] Base background uses `--colorNeutralBackground3`
- [ ] Shimmer gradient uses theme tokens:
  - `--colorNeutralBackground3` (start)
  - `--colorNeutralBackground2` (middle light)
  - `--colorNeutralBackground1` (middle bright)
  - `--colorNeutralBackground3` (end)

**Expected behavior:**
- Skeleton should be visible but subtle in all themes
- Shimmer animation should be smooth
- Should not be too bright or too dark in any theme

**How to test:**
1. Refresh the app to see initial loading state
2. Observe skeleton loaders for value selector
3. Verify shimmer effect is visible and smooth

---

### ✅ SuccessOverlay Component

**Elements to check:**
- [ ] Overlay background uses `colorNeutralBackgroundAlpha`
- [ ] Message container background uses `colorNeutralBackground1`
- [ ] Success icon uses `colorPaletteGreenForeground1`
- [ ] Message text uses `colorNeutralForeground1`

**Expected behavior:**
- Overlay should dim the background appropriately
- Message container should be clearly visible
- Success icon should be green and prominent
- Text should be readable

**How to test:**
1. Submit a card successfully
2. Observe the success overlay
3. Verify all colors are appropriate for the theme

---

### ✅ ErrorMessage Component

**Elements to check:**
- [ ] Uses Fluent UI MessageBar component (theme-aware by default)
- [ ] Error intent styling adapts to theme
- [ ] Retry and dismiss buttons use theme colors

**Expected behavior:**
- Error message should be clearly visible
- Red/error color should be appropriate for theme
- Buttons should be interactive and visible

**How to test:**
1. Disconnect network or cause an API error
2. Observe error message display
3. Test retry and dismiss buttons

---

### ✅ Animation Styles (animations.css)

**Elements to check:**
- [ ] Shimmer effect uses theme tokens:
  - `--colorNeutralBackground1Hover`
  - `--colorNeutralBackground1Pressed`
- [ ] Hover lift shadow uses `--shadow16`
- [ ] Focus ring uses `--colorStrokeFocus2`

**Expected behavior:**
- All animations should be visible in all themes
- Focus indicators should be clearly visible
- Hover effects should provide appropriate feedback

**How to test:**
1. Tab through interactive elements (keyboard navigation)
2. Verify focus rings are visible
3. Hover over buttons and cards
4. Verify hover effects are appropriate

---

## Verification Checklist

### Light Theme ☀️
- [ ] All components render correctly
- [ ] Text is readable (dark on light)
- [ ] Colors are vibrant and appropriate
- [ ] Animations are smooth
- [ ] No visual glitches

### Dark Theme 🌙
- [ ] All components render correctly
- [ ] Text is readable (light on dark)
- [ ] Colors are adjusted for dark background
- [ ] Animations are smooth
- [ ] No visual glitches
- [ ] No "flashbang" bright elements

### High Contrast Theme ⚡
- [ ] All components render correctly
- [ ] Maximum contrast between elements
- [ ] Text is highly readable
- [ ] Interactive elements are clearly distinguishable
- [ ] Animations respect reduced motion if enabled
- [ ] No elements are invisible or hard to see

## Common Issues to Watch For

### ❌ Hardcoded Colors
- White backgrounds in dark theme
- Black text in dark theme
- Colors that don't change with theme

### ❌ Insufficient Contrast
- Text that's hard to read
- Buttons that blend into background
- Disabled states that look enabled

### ❌ Animation Issues
- Animations that are too bright/dark
- Shimmer effects that are invisible
- Focus indicators that don't show

## Reporting Issues

If you find any theme adaptation issues:

1. **Document the issue:**
   - Which component?
   - Which theme mode?
   - What's wrong?
   - Screenshot if possible

2. **Check the code:**
   - Is a hardcoded color being used?
   - Is the correct theme token being used?
   - Is the CSS variable defined?

3. **Fix the issue:**
   - Replace hardcoded colors with theme tokens
   - Use Fluent UI tokens from `@fluentui/react-components`
   - Test in all three theme modes

## Success Criteria

✅ All components use theme tokens (no hardcoded colors)
✅ All components are readable in all three themes
✅ Animations adapt to theme colors
✅ Confetti uses theme-appropriate colors
✅ Focus indicators are visible in all themes
✅ No visual glitches or broken layouts

## Notes

- Theme tokens are CSS variables provided by Fluent UI's `FluentProvider`
- The `useThemeColors` hook provides type-safe access to theme colors
- All theme tokens start with `--color` prefix in CSS
- Fluent UI components automatically adapt to theme changes
