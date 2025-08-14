import React from 'react';
import ReactDOM from 'react-dom/client';
import { DashboardPages } from './src/presentation/pages/dashboard/dashboard';
import "./src/presentation/style/global.css"
import { AuthProvider } from './src/presentation/pages/middleware';
import './i18n';

ReactDOM.createRoot(document.getElementById("root")).render(
  <AuthProvider>
    <DashboardPages />
  </AuthProvider>
);
