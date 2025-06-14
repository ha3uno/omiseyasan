import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './index.css';
import App from './App';
import HistoryPage from './HistoryPage';
import ProductDetailPage from './ProductDetailPage';
import UserRegisterPage from './UserRegisterPage';
import CartPage from './CartPage';
import OrderConfirmationPage from './OrderConfirmationPage';
import PaymentCompletePage from './PaymentCompletePage';
import OrderHistoryPage from './OrderHistoryPage';
import { CartProvider } from './CartContext';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <CartProvider>
      <Router>
        <Routes>
          <Route path="/" element={<App />} />
          <Route path="/history" element={<HistoryPage />} />
          <Route path="/products/:id" element={<ProductDetailPage />} />
          <Route path="/register" element={<UserRegisterPage />} />
          <Route path="/cart" element={<CartPage />} />
          <Route path="/order" element={<OrderConfirmationPage />} />
          <Route path="/payment-complete" element={<PaymentCompletePage />} />
          <Route path="/order-history" element={<OrderHistoryPage />} />
        </Routes>
      </Router>
    </CartProvider>
  </React.StrictMode>
);