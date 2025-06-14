import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import './index.css';
import App from './App';
import HistoryPage from './HistoryPage';
import ProductDetailPage from './ProductDetailPage';
import UserRegisterPage from './UserRegisterPage';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <Router>
      <Routes>
        <Route path="/" element={<App />} />
        <Route path="/history" element={<HistoryPage />} />
        <Route path="/products/:id" element={<ProductDetailPage />} />
        <Route path="/register" element={<UserRegisterPage />} />
      </Routes>
    </Router>
  </React.StrictMode>
);