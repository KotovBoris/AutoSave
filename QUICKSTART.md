# Quick Start Guide

## Installation

1. Install dependencies:
```bash
npm install
```

2. Run the development server:
```bash
npm run dev
```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Project Structure

```
vtb/
├── app/                      # Next.js app directory
│   ├── (auth)/              # Authentication pages (register, login)
│   ├── (dashboard)/         # Dashboard pages (protected routes)
│   │   ├── dashboard/       # Main dashboard
│   │   ├── goals/           # Goals management
│   │   ├── loans/           # Loans management
│   │   ├── accounts/        # Accounts page
│   │   ├── history/         # Operations history
│   │   └── settings/        # Settings page
│   ├── onboarding/          # Onboarding flow
│   │   ├── banks/           # Bank connection
│   │   ├── analysis/        # Transaction analysis
│   │   ├── salaries/        # Salary confirmation
│   │   └── result/          # Analysis result
│   └── layout.tsx           # Root layout
├── components/              # React components
│   ├── layout/             # Layout components
│   │   ├── Topbar.tsx      # Top navigation bar
│   │   ├── Sidebar.tsx     # Side navigation
│   │   └── NotificationsDropdown.tsx
│   ├── ui/                 # UI components
│   │   └── Modal.tsx       # Modal component
│   └── modals/             # Modal components
│       └── EmergencyWithdrawModal.tsx
├── lib/                    # Utilities
│   ├── api.ts              # Mock API functions
│   ├── mockData.ts         # Mock data
│   └── utils.ts            # Utility functions
└── types/                  # TypeScript types
    └── index.ts            # Type definitions
```

## Features

### ✅ Implemented

- **Authentication**: Registration and login pages
- **Onboarding Flow**: 
  - Bank connection (mock OAuth)
  - Transaction analysis simulation
  - Salary confirmation
  - Savings capacity calculation
- **Dashboard**: 
  - Statistics cards (total balance, goals, loans)
  - Active goal display
  - Next actions
- **Goals Management**: 
  - Create, edit, delete goals
  - Reorder goals
  - Progress tracking
  - Calculation of savings plan
- **Loans Management**: 
  - Add, edit, delete loans
  - Auto-payment settings
  - Loan schedule calculation
- **Accounts**: View all accounts and transactions
- **History**: View operation history
- **Settings**: User profile and bank connections
- **Emergency Withdrawal**: Modal for emergency fund withdrawal
- **Notifications**: Dropdown with recent operations

## Mock Data

The application uses mock data for demonstration. All API calls are simulated with delays to mimic real API behavior. To connect to a real backend:

1. Update `lib/api.ts` with real API endpoints
2. Replace mock data in `lib/mockData.ts` with real data fetching
3. Update authentication to use real auth tokens

## API Structure

The mock API follows this structure:

```typescript
// Auth
authAPI.register(email, password)
authAPI.login(email, password)
authAPI.getCurrentUser()

// Banks
banksAPI.getBanks()
banksAPI.connectBank(bankId)
banksAPI.disconnectBank(bankId)

// Accounts
accountsAPI.getAccounts()
accountsAPI.getAccount(id)
accountsAPI.getAccountTransactions(accountId)

// Goals
goalsAPI.getGoals()
goalsAPI.createGoal(goal)
goalsAPI.updateGoal(id, updates)
goalsAPI.deleteGoal(id)
goalsAPI.reorderGoals(goalIds)

// Loans
loansAPI.getLoans()
loansAPI.createLoan(loan)
loansAPI.updateLoan(id, updates)
loansAPI.deleteLoan(id)

// Operations
operationsAPI.getOperations()
operationsAPI.emergencyWithdraw(amount)

// Analysis
analysisAPI.getSalaries()
analysisAPI.getSavingsCapacity()
```

## Styling

The application uses Tailwind CSS for styling. Custom colors are defined in `tailwind.config.ts`:

- Primary: Indigo/Blue colors for main actions
- Danger: Red colors for destructive actions
- Custom utility classes in `app/globals.css`

## Next Steps

1. **Connect to Real Backend**: Replace mock API with real endpoints
2. **Add Charts**: Implement Recharts for goal progress and loan schedules
3. **Add Authentication**: Implement real JWT authentication
4. **Add OAuth**: Implement real bank OAuth flows
5. **Add Error Handling**: Improve error handling and user feedback
6. **Add Tests**: Add unit and integration tests
7. **Add Mobile Support**: Make the application responsive for mobile devices

## Troubleshooting

### TypeScript Errors
If you see TypeScript errors about missing modules, make sure to run `npm install` first.

### Build Errors
If you encounter build errors, try:
```bash
rm -rf .next node_modules
npm install
npm run dev
```

### Port Already in Use
If port 3000 is already in use, you can change it by modifying the `dev` script in `package.json`:
```json
"dev": "next dev -p 3001"
```

