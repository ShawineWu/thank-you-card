/**
 * ConfettiEffect Component Usage Examples
 * 
 * This file demonstrates various ways to use the ConfettiEffect component
 * in the Teams App for celebration animations.
 */

import { useState } from 'react';
import { Button } from '@fluentui/react-components';
import { ConfettiEffect } from './ConfettiEffect';

/**
 * Example 1: Basic Usage
 * 
 * Simple confetti effect with default configuration
 */
export function BasicConfettiExample() {
  const [showConfetti, setShowConfetti] = useState(false);

  return (
    <div>
      <Button 
        appearance="primary"
        onClick={() => setShowConfetti(true)}
      >
        Celebrate! 🎉
      </Button>
      
      <ConfettiEffect 
        active={showConfetti}
        onComplete={() => setShowConfetti(false)}
      />
    </div>
  );
}

/**
 * Example 2: Custom Configuration
 * 
 * Confetti with custom particle count, spread, and duration
 */
export function CustomConfettiExample() {
  const [showConfetti, setShowConfetti] = useState(false);

  return (
    <div>
      <Button 
        appearance="primary"
        onClick={() => setShowConfetti(true)}
      >
        Big Celebration! 🎊
      </Button>
      
      <ConfettiEffect 
        active={showConfetti}
        onComplete={() => setShowConfetti(false)}
        config={{
          particleCount: 200,
          spread: 90,
          duration: 4000,
        }}
      />
    </div>
  );
}

/**
 * Example 3: Custom Colors
 * 
 * Confetti with specific brand colors
 */
export function CustomColorConfettiExample() {
  const [showConfetti, setShowConfetti] = useState(false);

  return (
    <div>
      <Button 
        appearance="primary"
        onClick={() => setShowConfetti(true)}
      >
        Brand Colors! 🌈
      </Button>
      
      <ConfettiEffect 
        active={showConfetti}
        onComplete={() => setShowConfetti(false)}
        config={{
          colors: ['#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A', '#98D8C8'],
        }}
      />
    </div>
  );
}

/**
 * Example 4: Custom Origin Point
 * 
 * Confetti originating from a specific location
 */
export function CustomOriginConfettiExample() {
  const [showConfetti, setShowConfetti] = useState(false);

  return (
    <div>
      <Button 
        appearance="primary"
        onClick={() => setShowConfetti(true)}
      >
        From Bottom! 🚀
      </Button>
      
      <ConfettiEffect 
        active={showConfetti}
        onComplete={() => setShowConfetti(false)}
        config={{
          origin: { x: 0.5, y: 0.9 }, // Bottom center
          spread: 120,
        }}
      />
    </div>
  );
}

/**
 * Example 5: Integration with Form Success
 * 
 * Typical usage in a form submission success scenario
 */
export function FormSuccessConfettiExample() {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);

  const handleSubmit = async () => {
    setIsSubmitting(true);
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    setIsSubmitting(false);
    setShowSuccess(true);
  };

  return (
    <div>
      <Button 
        appearance="primary"
        disabled={isSubmitting}
        onClick={handleSubmit}
      >
        {isSubmitting ? 'Submitting...' : 'Submit Form'}
      </Button>
      
      {showSuccess && (
        <>
          <div style={{ marginTop: '20px', color: 'green' }}>
            ✓ Success! Your form has been submitted.
          </div>
          
          <ConfettiEffect 
            active={showSuccess}
            onComplete={() => setShowSuccess(false)}
          />
        </>
      )}
    </div>
  );
}

/**
 * Example 6: Theme-Aware Confetti
 * 
 * Confetti that automatically uses theme colors (default behavior)
 * The component will extract colors from the current Fluent UI theme
 */
export function ThemeAwareConfettiExample() {
  const [showConfetti, setShowConfetti] = useState(false);

  return (
    <div>
      <Button 
        appearance="primary"
        onClick={() => setShowConfetti(true)}
      >
        Theme Colors! 🎨
      </Button>
      
      <p style={{ fontSize: '12px', marginTop: '8px' }}>
        Confetti colors will match your current Teams theme
        (light, dark, or high contrast)
      </p>
      
      <ConfettiEffect 
        active={showConfetti}
        onComplete={() => setShowConfetti(false)}
        // No colors specified - will use theme colors automatically
      />
    </div>
  );
}
