import React, { useState, useEffect } from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './index.css';
import App from './App';
import HistoryPage from './HistoryPage';
import UpdateHistoryListPage from './UpdateHistoryListPage';
import UpdateHistoryRegisterPage from './UpdateHistoryRegisterPage';
import ProductDetailPage from './ProductDetailPage';
import UserRegisterPage from './UserRegisterPage';
import CartPage from './CartPage';
import OrderConfirmationPage from './OrderConfirmationPage';
import PaymentCompletePage from './PaymentCompletePage';
import OrderHistoryPage from './OrderHistoryPage';
import LoginPage from './LoginPage';
import MyAccountPage from './MyAccountPage';
import { CartProvider } from './CartContext';

interface User {
  id: number;
  name: string;
  address: string;
  phoneNumber: string;
  email: string;
}

const USER_STORAGE_KEY = 'omiseyasan-user';

function AppWithAuth() {
  const [loggedInUser, setLoggedInUser] = useState<User | null>(null);

  // Load user from localStorage on component mount
  useEffect(() => {
    try {
      const savedUser = localStorage.getItem(USER_STORAGE_KEY);
      if (savedUser) {
        const parsedUser = JSON.parse(savedUser);
        setLoggedInUser(parsedUser);
      }
    } catch (error) {
      console.warn('Failed to load user from localStorage:', error);
      localStorage.removeItem(USER_STORAGE_KEY);
    }
  }, []);

  // Save user to localStorage whenever loggedInUser changes
  useEffect(() => {
    try {
      if (loggedInUser) {
        localStorage.setItem(USER_STORAGE_KEY, JSON.stringify(loggedInUser));
      } else {
        localStorage.removeItem(USER_STORAGE_KEY);
      }
    } catch (error) {
      console.warn('Failed to save user to localStorage:', error);
    }
  }, [loggedInUser]);

  const handleLogin = (user: User) => {
    setLoggedInUser(user);
  };

  const handleLogout = () => {
    setLoggedInUser(null);
  };

  return (
    <Routes>
      <Route path="/" element={<App loggedInUser={loggedInUser} onLogin={handleLogin} onLogout={handleLogout} />} />
      <Route path="/history" element={<HistoryPage />} />
      <Route path="/update-history" element={<UpdateHistoryListPage />} />
      <Route path="/update-history/register" element={<UpdateHistoryRegisterPage />} />
      <Route path="/products/:id" element={<ProductDetailPage />} />
      <Route path="/register" element={<UserRegisterPage />} />
      <Route path="/cart" element={<CartPage />} />
      <Route path="/order" element={<OrderConfirmationPage />} />
      <Route path="/payment-complete" element={<PaymentCompletePage />} />
      <Route path="/order-history" element={<OrderHistoryPage />} />
      <Route path="/login" element={<LoginPage onLogin={handleLogin} />} />
      <Route path="/my-account" element={<MyAccountPage loggedInUser={loggedInUser} onLogout={handleLogout} />} />
    </Routes>
  );
}

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <CartProvider>
      <Router>
        <AppWithAuth />
      </Router>
    </CartProvider>
  </React.StrictMode>
);