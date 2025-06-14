import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import './App.css';

interface UserRegistrationForm {
  name: string;
  address: string;
  phoneNumber: string;
  email: string;
}

function UserRegisterPage() {
  const [formData, setFormData] = useState<UserRegistrationForm>({
    name: '',
    address: '',
    phoneNumber: '',
    email: ''
  });
  
  const [submitting, setSubmitting] = useState<boolean>(false);
  const [success, setSuccess] = useState<boolean>(false);
  const [error, setError] = useState<string>('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // Basic validation
    if (!formData.name.trim() || !formData.email.trim()) {
      setError('お名前とメールアドレスは必須項目です。');
      return;
    }

    try {
      setSubmitting(true);
      setError('');
      setSuccess(false);

      const response = await fetch('/api/users/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || `HTTP error! status: ${response.status}`);
      }

      const registeredUser = await response.json();
      console.log('User registered:', registeredUser);
      
      // Reset form and show success
      setFormData({
        name: '',
        address: '',
        phoneNumber: '',
        email: ''
      });
      setSuccess(true);

    } catch (err) {
      setError(err instanceof Error ? err.message : '登録エラーが発生しました');
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
          <Link to="/cart" className="nav-link">🛒 カートを見る</Link>
          <Link to="/history" className="nav-link">📝 更新履歴</Link>
        </nav>
      </header>
      
      <main>
        <div className="register-container">
          <h2>アカウント登録</h2>
          <p>おみせやさん♪でお買い物をするためのアカウントを作成してください。</p>
          
          {success && (
            <div className="success-message">
              <p>✅ アカウントが正常に登録されました！</p>
              <Link to="/" className="nav-link">商品一覧を見る</Link>
            </div>
          )}
          
          {!success && (
            <form onSubmit={handleSubmit} className="register-form">
              <div className="form-group">
                <label htmlFor="name">お名前 *</label>
                <input
                  type="text"
                  id="name"
                  name="name"
                  value={formData.name}
                  onChange={handleInputChange}
                  placeholder="山田太郎"
                  required
                />
              </div>

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
                <label htmlFor="phoneNumber">電話番号</label>
                <input
                  type="tel"
                  id="phoneNumber"
                  name="phoneNumber"
                  value={formData.phoneNumber}
                  onChange={handleInputChange}
                  placeholder="090-1234-5678"
                />
              </div>

              <div className="form-group">
                <label htmlFor="address">住所</label>
                <textarea
                  id="address"
                  name="address"
                  value={formData.address}
                  onChange={handleInputChange}
                  placeholder="〒123-4567 東京都渋谷区..."
                  rows={3}
                />
              </div>

              {error && <p className="form-error">{error}</p>}

              <button 
                type="submit" 
                disabled={submitting}
                className="submit-button"
              >
                {submitting ? '登録中...' : '🚀 アカウントを登録'}
              </button>
            </form>
          )}
          
          <div className="navigation-links">
            <Link to="/" className="nav-link">← ホームに戻る</Link>
          </div>
        </div>
      </main>
    </div>
  );
}

export default UserRegisterPage;