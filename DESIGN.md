---
name: Terminal Precision (High Contrast)
colors:
  surface: '#051424'
  surface-dim: '#051424'
  surface-bright: '#2c3a4c'
  surface-container-lowest: '#010f1f'
  surface-container-low: '#0d1c2d'
  surface-container: '#142538'
  surface-container-high: '#1c2f45'
  surface-container-highest: '#263a52'
  on-surface: '#e2e8f0'
  on-surface-variant: '#94a3b8'
  outline: '#475569'
  outline-variant: '#1e293b'
  primary: '#22d3ee'
  on-primary: '#083344'
  primary-container: '#0e7490'
  on-primary-container: '#cffafe'
  secondary: '#334155'
  on-secondary: '#f1f5f9'
  secondary-container: '#1e293b'
  on-secondary-container: '#cbd5e1'
  error: '#f87171'
  on-error: '#450a0a'
  success: '#4ade80'
  on-success: '#052e16'
  warning: '#fbbf24'
  on-warning: '#451a03'

typography:
  font-family: 'Geist, JetBrains Mono, monospace'
  base-size: 14px
  scale: 1.2

roundness: 4px (ROUND_FOUR)

principles:
  - High Contrast: Prioritize legibility with deep backgrounds and vibrant foregrounds.
  - Information Density: Maximize data visibility for technical workflows.
  - Status-at-a-glance: Use functional color (Cyan, Red, Green) for actionable state.
  - Minimalist Flatness: Use surface elevation shifts rather than shadows for depth.

components:
  TopAppBar:
    background: surface
    border-bottom: outline-variant
    text: primary
  BottomNavBar:
    background: surface-container
    active-item: secondary-container
    active-text: primary
---

# Terminal Precision Design System

## Brand Personality
Precise, technical, and engineering-focused. The aesthetic is inspired by modern terminal emulators and IDEs, optimized for high-density mobile information display.

## Visual Language
The system uses a dark-mode first palette centered on a deep navy slate (`#051424`) to provide maximum stability for the high-intensity Cyan (`#22d3ee`) accents.

## Typography
A mix of **Geist** for clean, modern legibility and **JetBrains Mono** for technical data. The type scale is optimized for high information density without sacrificing readability.

## Elevation & Depth
Elevation is expressed through background color shifts (`surface-container-low` to `highest`) rather than drop shadows, maintaining a flat, professional engineering aesthetic.
