# Terminal Precision Design System

## Brand Personality
- **Voice**: Precise, technical, engineering-focused.
- **Vibe**: High-density, terminal-inspired, professional, dark-mode centric.

## Visual Language
- **Color Palette**: 
  - **Surface**: `#051424` (Deep navy slate)
  - **Primary**: `#22d3ee` (Cyan highlight for interactivity and status)
  - **Secondary**: Darker slates and outline variants for structural separation.
- **Typography**: 
  - **Primary Font**: Geist / JetBrains Mono (Monospace and clean sans-serif mix)
  - **Scale**: Optimized for high-density mobile information.
- **Roundness**: `ROUND_FOUR` (Subtle 4px rounding for a modern yet rigid feel).
- **Style**: Flat, minimalist, no shadows. Elevation is handled via surface color shifts (container-low, container-high) rather than depth effects.

## Shared Components

### TopAppBar
- **Style**: Small, docked to top, full-width.
- **Visuals**: `#131313` background, border-bottom separator, cyan primary text.
- **Features**: Leading terminal icon, trailing status indicator ("CONNECTED").

### BottomNavBar
- **Style**: Label + Icon, docked to bottom, full-width.
- **Destinations**: PRs (merge icon), Workspaces (grid icon), Chat (forum icon), Settings (settings icon).
- **Interaction**: Scale-95 transition on active tap; secondary-container background for the active state.

## Design Principles
1. **Information Density**: Prioritize data over whitespace, but maintain legibility through clear borders and type scale.
2. **Status at a Glance**: Use the primary color (Cyan) and secondary status colors (Yellow for changes requested, Green for approved) sparingly to draw attention to actionable items.
3. **Consistency**: Use the established grid and spacing tokens to ensure all cards and lists align perfectly.

