import React from 'react';
import '../styles/animations.css';
import './SkeletonLoader.css';

export interface SkeletonLoaderProps {
  variant: 'text' | 'circular' | 'rectangular';
  width?: string | number;
  height?: string | number;
  animation?: 'pulse' | 'wave' | 'none';
}

/**
 * SkeletonLoader Component
 * 
 * Displays a loading placeholder with shimmer animation effect.
 * Supports multiple variants for different content types.
 * 
 * Requirements: 1.5, 5.1, 5.2
 */
export const SkeletonLoader: React.FC<SkeletonLoaderProps> = ({
  variant,
  width,
  height,
  animation = 'wave',
}) => {
  const getVariantStyles = (): React.CSSProperties => {
    const baseStyles: React.CSSProperties = {
      width: typeof width === 'number' ? `${width}px` : width,
      height: typeof height === 'number' ? `${height}px` : height,
    };

    switch (variant) {
      case 'text':
        return {
          ...baseStyles,
          width: width || '100%',
          height: height || '1em',
          borderRadius: '4px',
        };
      case 'circular':
        return {
          ...baseStyles,
          width: width || '40px',
          height: height || '40px',
          borderRadius: '50%',
        };
      case 'rectangular':
        return {
          ...baseStyles,
          width: width || '100%',
          height: height || '100px',
          borderRadius: '8px',
        };
      default:
        return baseStyles;
    }
  };

  const getAnimationClass = (): string => {
    switch (animation) {
      case 'wave':
        return 'skeleton-wave';
      case 'pulse':
        return 'skeleton-pulse';
      case 'none':
        return '';
      default:
        return 'skeleton-wave';
    }
  };

  return (
    <div
      className={`skeleton-loader ${getAnimationClass()}`}
      style={getVariantStyles()}
      role="status"
      aria-label="Loading..."
      aria-busy="true"
    />
  );
};
