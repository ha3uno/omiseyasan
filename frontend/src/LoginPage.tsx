import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from './CartContext';
import './App.css';

interface LoginForm {
  email: string;
  password: string;
}

interface LoginPageProps {
  onLogin: (user: any) => void;
}

function LoginPage({ onLogin }: LoginPageProps) {
  const navigate = useNavigate();
  const { getTotalQuantity } = useCart();
  
  const [formData, setFormData] = useState<LoginForm>({
    email: '',
    password: ''
  });
  
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    if (!formData.email.trim() || !formData.password.trim()) {
      setError('メールアドレスとパスワードを入力してください。');
      return;
    }

    try {
      setSubmitting(true);
      setError('');

      const response = await fetch('/api/users/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        if (response.status === 401) {
          throw new Error('メールアドレスまたはパスワードが間違っています。');
        }
        const errorText = await response.text();
        throw new Error(errorText || `HTTP error! status: ${response.status}`);
      }

      const user = await response.json();
      console.log('Login successful:', user);
      
      // Call the login callback to update global state
      onLogin(user);
      
      // Navigate to my account page
      navigate('/my-account');

    } catch (err) {
      setError(err instanceof Error ? err.message : 'ログインエラーが発生しました');
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>おみせやさん♪</h1>
        <nav className="header-nav">
          <Link to="/" className="nav-link">🏠 ホーム</Link>
          <Link to="/cart" className="nav-link">🛒 カート ({getTotalQuantity()})</Link>
          <Link to="/register" className="nav-link">👤 アカウント登録</Link>
          <Link to="/history" className="nav-link">📝 更新履歴</Link>
        </nav>
      </header>
      
      <main>
        <div className="login-container">
          <h2>ログイン</h2>
          <p>アカウントにログインして、お買い物をお楽しみください。</p>
          
          <form onSubmit={handleSubmit} className="login-form">
            <div className="form-group">
              <label htmlFor="email">メールアドレス *</label>
              <input
                type="email"
                id="email"
                name="email"
                value={formData.email}
                onChange={handleInputChange}
                placeholder="example@email.com"
                required
              />
            </div>

            <div className="form-group">
              <label htmlFor="password">パスワード *</label>
              <input
                type="password"
                id="password"
                name="password"
                value={formData.password}
                onChange={handleInputChange}
                placeholder="パスワードを入力してください"
                required
              />
            </div>

            {error && <p className="form-error">{error}</p>}

            <button 
              type="submit" 
              disabled={submitting}
              className="submit-button"
            >
              {submitting ? 'ログイン中...' : '🔑 ログイン'}
            </button>
          </form>
          
          <div className="login-links">
            <p>アカウントをお持ちでない方は</p>
            <Link to="/register" className="nav-link">新規アカウント登録</Link>
          </div>
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← ホームに戻る</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

export default LoginPage;