# Fitness Hack Frontend

A modern React TypeScript frontend for the Fitness Hack application, built with MobX for state management and Tailwind CSS for styling.

## Tech Stack

- **React 18** with TypeScript
- **MobX** for state management
- **React Router** for navigation
- **Tailwind CSS** for styling
- **Vite** for build tooling
- **Axios** for API calls
- **Lucide React** for icons

## Project Structure

```
frontend/
├── src/
│   ├── components/     # Reusable UI components
│   ├── pages/         # Page components
│   ├── stores/        # MobX stores
│   ├── services/      # API services
│   ├── types/         # TypeScript type definitions
│   ├── utils/         # Utility functions
│   ├── App.tsx        # Main app component
│   ├── main.tsx       # App entry point
│   └── index.css      # Global styles
├── package.json       # Dependencies and scripts
├── vite.config.ts     # Vite configuration
├── tailwind.config.js # Tailwind CSS configuration
└── tsconfig.json      # TypeScript configuration
```

## Getting Started

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Open your browser and navigate to `http://localhost:3000`

### Build for Production

```bash
npm run build
```

The built files will be in the `dist/` directory.

## Development

### Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

### State Management

The app uses MobX for state management with the following stores:

- **UserStore** - Manages user authentication and profile
- **WorkoutStore** - Manages workout data
- **ExerciseStore** - Manages exercise data
- **RootStore** - Combines all stores

### API Integration

The frontend is configured to proxy API calls to the backend at `http://localhost:8080`. API calls are made through services in the `src/services/` directory.

### Styling

The app uses Tailwind CSS with custom components defined in `src/index.css`. Custom utility classes are available:

- `.btn-primary` - Primary button styling
- `.btn-secondary` - Secondary button styling
- `.card` - Card container styling
- `.input-field` - Form input styling

## Features

- **Dashboard** - Overview of fitness data
- **Workouts** - Manage workout routines
- **Exercises** - Browse and manage exercises
- **Profile** - User account management
- **Responsive Design** - Works on desktop and mobile
- **Type Safety** - Full TypeScript support

## Contributing

1. Follow the existing code structure
2. Use TypeScript for all new code
3. Use MobX for state management
4. Follow the component naming conventions
5. Add proper TypeScript types for all props and state 