// Unified color configuration for company values and credos
// Each value has a unique color for consistent display across the site

export interface ValueColorConfig {
  bg: string;       // Background color for badges/tags (Tailwind class)
  text: string;     // Text color (Tailwind class)
  border: string;   // Border color (Tailwind class)
  bgHover: string;  // Hover background color (Tailwind class)
  hex: string;      // Hex color for charts
  // CSS values for inline styles
  bgCss: string;    // Background color CSS value
  textCss: string;  // Text color CSS value
  borderCss: string; // Border color CSS value
}

// 5 Values
const VALUE_COLORS: Record<string, ValueColorConfig> = {
  'Make an Impact': {
    bg: 'bg-violet-100',
    text: 'text-violet-700',
    border: 'border-violet-200',
    bgHover: 'hover:bg-violet-200',
    hex: '#8b5cf6',
    bgCss: '#ede9fe',
    textCss: '#6d28d9',
    borderCss: '#ddd6fe',
  },
  'Strive for Excellence': {
    bg: 'bg-amber-100',
    text: 'text-amber-700',
    border: 'border-amber-200',
    bgHover: 'hover:bg-amber-200',
    hex: '#f59e0b',
    bgCss: '#fef3c7',
    textCss: '#b45309',
    borderCss: '#fde68a',
  },
  'Stand Together': {
    bg: 'bg-emerald-100',
    text: 'text-emerald-700',
    border: 'border-emerald-200',
    bgHover: 'hover:bg-emerald-200',
    hex: '#10b981',
    bgCss: '#d1fae5',
    textCss: '#047857',
    borderCss: '#a7f3d0',
  },
  'Be Open-Minded': {
    bg: 'bg-sky-100',
    text: 'text-sky-700',
    border: 'border-sky-200',
    bgHover: 'hover:bg-sky-200',
    hex: '#0ea5e9',
    bgCss: '#e0f2fe',
    textCss: '#0369a1',
    borderCss: '#bae6fd',
  },
  'Stay Grounded': {
    bg: 'bg-rose-100',
    text: 'text-rose-700',
    border: 'border-rose-200',
    bgHover: 'hover:bg-rose-200',
    hex: '#f43f5e',
    bgCss: '#ffe4e6',
    textCss: '#be123c',
    borderCss: '#fecdd3',
  },
};

// 10 Credos
const CREDO_COLORS: Record<string, ValueColorConfig> = {
  'Bias for Action': {
    bg: 'bg-orange-100',
    text: 'text-orange-700',
    border: 'border-orange-200',
    bgHover: 'hover:bg-orange-200',
    hex: '#f97316',
    bgCss: '#ffedd5',
    textCss: '#c2410c',
    borderCss: '#fed7aa',
  },
  'Customer Centric': {
    bg: 'bg-pink-100',
    text: 'text-pink-700',
    border: 'border-pink-200',
    bgHover: 'hover:bg-pink-200',
    hex: '#ec4899',
    bgCss: '#fce7f3',
    textCss: '#be185d',
    borderCss: '#fbcfe8',
  },
  'Think Strategically': {
    bg: 'bg-indigo-100',
    text: 'text-indigo-700',
    border: 'border-indigo-200',
    bgHover: 'hover:bg-indigo-200',
    hex: '#6366f1',
    bgCss: '#e0e7ff',
    textCss: '#4338ca',
    borderCss: '#c7d2fe',
  },
  'Deep Dive': {
    bg: 'bg-cyan-100',
    text: 'text-cyan-700',
    border: 'border-cyan-200',
    bgHover: 'hover:bg-cyan-200',
    hex: '#06b6d4',
    bgCss: '#cffafe',
    textCss: '#0e7490',
    borderCss: '#a5f3fc',
  },
  'Invent and Simplify': {
    bg: 'bg-teal-100',
    text: 'text-teal-700',
    border: 'border-teal-200',
    bgHover: 'hover:bg-teal-200',
    hex: '#14b8a6',
    bgCss: '#ccfbf1',
    textCss: '#0f766e',
    borderCss: '#99f6e4',
  },
  'Earn Trust': {
    bg: 'bg-blue-100',
    text: 'text-blue-700',
    border: 'border-blue-200',
    bgHover: 'hover:bg-blue-200',
    hex: '#3b82f6',
    bgCss: '#dbeafe',
    textCss: '#1d4ed8',
    borderCss: '#bfdbfe',
  },
  'Take Ownership': {
    bg: 'bg-purple-100',
    text: 'text-purple-700',
    border: 'border-purple-200',
    bgHover: 'hover:bg-purple-200',
    hex: '#a855f7',
    bgCss: '#f3e8ff',
    textCss: '#7e22ce',
    borderCss: '#e9d5ff',
  },
  'Challenge Disagree and Commit': {
    bg: 'bg-red-100',
    text: 'text-red-700',
    border: 'border-red-200',
    bgHover: 'hover:bg-red-200',
    hex: '#ef4444',
    bgCss: '#fee2e2',
    textCss: '#b91c1c',
    borderCss: '#fecaca',
  },
  'Learn and Be Curious': {
    bg: 'bg-lime-100',
    text: 'text-lime-700',
    border: 'border-lime-200',
    bgHover: 'hover:bg-lime-200',
    hex: '#84cc16',
    bgCss: '#ecfccb',
    textCss: '#4d7c0f',
    borderCss: '#d9f99d',
  },
  'Do More with Less': {
    bg: 'bg-fuchsia-100',
    text: 'text-fuchsia-700',
    border: 'border-fuchsia-200',
    bgHover: 'hover:bg-fuchsia-200',
    hex: '#d946ef',
    bgCss: '#fae8ff',
    textCss: '#a21caf',
    borderCss: '#f5d0fe',
  },
};

// Combined colors map
export const VALUE_COLOR_MAP: Record<string, ValueColorConfig> = {
  ...VALUE_COLORS,
  ...CREDO_COLORS,
};

// Default color for unknown values
export const DEFAULT_VALUE_COLOR: ValueColorConfig = {
  bg: 'bg-gray-100',
  text: 'text-gray-700',
  border: 'border-gray-200',
  bgHover: 'hover:bg-gray-200',
  hex: '#6b7280',
  bgCss: '#f3f4f6',
  textCss: '#374151',
  borderCss: '#e5e7eb',
};

// Helper function to get color config by value name
export function getValueColor(valueName: string): ValueColorConfig {
  return VALUE_COLOR_MAP[valueName] || DEFAULT_VALUE_COLOR;
}

// Helper function to get combined class string for a badge
export function getValueBadgeClasses(valueName: string): string {
  const color = getValueColor(valueName);
  return `${color.bg} ${color.text} ${color.border} ${color.bgHover}`;
}

// Helper function to get hex color for charts
export function getValueHexColor(valueName: string): string {
  const color = getValueColor(valueName);
  return color.hex;
}
