/**
 * SkeletonLoader Usage Examples
 * 
 * This file demonstrates various use cases for the SkeletonLoader component.
 * These examples can be used in the actual components when implementing loading states.
 */

import React from 'react';
import { SkeletonLoader } from './SkeletonLoader';

/**
 * Example 1: Text Loading Skeleton
 * Use for loading text content like card messages, descriptions, etc.
 */
export const TextSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px' }}>
    <SkeletonLoader variant="text" width="100%" height="1em" />
    <SkeletonLoader variant="text" width="90%" height="1em" />
    <SkeletonLoader variant="text" width="95%" height="1em" />
  </div>
);

/**
 * Example 2: Circular Avatar Skeleton
 * Use for loading user avatars in recipient selector or card history
 */
export const CircularSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px', display: 'flex', gap: '12px' }}>
    <SkeletonLoader variant="circular" width={40} height={40} />
    <SkeletonLoader variant="circular" width={40} height={40} />
    <SkeletonLoader variant="circular" width={40} height={40} />
  </div>
);

/**
 * Example 3: Rectangular Card Skeleton
 * Use for loading card previews or history cards
 */
export const RectangularSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px' }}>
    <SkeletonLoader variant="rectangular" width="100%" height={200} />
  </div>
);

/**
 * Example 4: Complete Card Form Loading State
 * Shows how to use multiple skeletons to create a complete loading state
 */
export const CardFormSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px', maxWidth: '600px' }}>
    {/* Header */}
    <SkeletonLoader variant="text" width="200px" height="24px" />
    
    <div style={{ marginTop: '20px' }}>
      {/* Recipient selector */}
      <SkeletonLoader variant="text" width="100px" height="16px" />
      <div style={{ marginTop: '8px', display: 'flex', gap: '8px' }}>
        <SkeletonLoader variant="circular" width={32} height={32} />
        <SkeletonLoader variant="circular" width={32} height={32} />
        <SkeletonLoader variant="circular" width={32} height={32} />
      </div>
    </div>
    
    <div style={{ marginTop: '20px' }}>
      {/* Message field */}
      <SkeletonLoader variant="text" width="120px" height="16px" />
      <div style={{ marginTop: '8px' }}>
        <SkeletonLoader variant="rectangular" width="100%" height={100} />
      </div>
    </div>
    
    <div style={{ marginTop: '20px' }}>
      {/* Value selector */}
      <SkeletonLoader variant="text" width="140px" height="16px" />
      <div style={{ marginTop: '8px', display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
        <SkeletonLoader variant="rectangular" width={80} height={32} />
        <SkeletonLoader variant="rectangular" width={100} height={32} />
        <SkeletonLoader variant="rectangular" width={90} height={32} />
      </div>
    </div>
    
    <div style={{ marginTop: '20px' }}>
      {/* Preview */}
      <SkeletonLoader variant="text" width="100px" height="16px" />
      <div style={{ marginTop: '8px' }}>
        <SkeletonLoader variant="rectangular" width="100%" height={250} />
      </div>
    </div>
  </div>
);

/**
 * Example 5: History Grid Loading State
 * Shows skeleton for card history grid
 */
export const HistoryGridSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px', display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(250px, 1fr))', gap: '16px' }}>
    <SkeletonLoader variant="rectangular" width="100%" height={180} />
    <SkeletonLoader variant="rectangular" width="100%" height={180} />
    <SkeletonLoader variant="rectangular" width="100%" height={180} />
    <SkeletonLoader variant="rectangular" width="100%" height={180} />
  </div>
);

/**
 * Example 6: Pulse Animation Variant
 * Alternative animation style using pulse instead of wave
 */
export const PulseSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px' }}>
    <SkeletonLoader variant="text" width="100%" animation="pulse" />
    <SkeletonLoader variant="text" width="90%" animation="pulse" />
    <SkeletonLoader variant="rectangular" width="100%" height={150} animation="pulse" />
  </div>
);

/**
 * Example 7: No Animation (for reduced motion preference)
 * Shows skeleton without animation
 */
export const NoAnimationSkeletonExample: React.FC = () => (
  <div style={{ padding: '20px' }}>
    <SkeletonLoader variant="text" width="100%" animation="none" />
    <SkeletonLoader variant="rectangular" width="100%" height={150} animation="none" />
  </div>
);
