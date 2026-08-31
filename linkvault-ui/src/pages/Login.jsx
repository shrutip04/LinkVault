import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '../api';
import styles from './Auth.module.css';

export default function Login() {
  const [form, setForm] = useState({ email: '', password: '' });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await api.post('/login', form);
      localStorage.setItem('token', res.data.token);
      localStorage.setItem('username', res.data.username);
      toast.success(`Welcome back, ${res.data.username}!`);
      navigate('/');
    } catch (err) {
      toast.error(err.response?.data?.error || 'Login failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <h1 className={styles.logo}>🔗 LinkVault</h1>
        <h2 className={styles.title}>Welcome back</h2>
        <form onSubmit={handleSubmit} className={styles.form}>
          <input
            className={styles.input}
            type="email"
            placeholder="Email"
            value={form.email}
            onChange={e => setForm({...form, email: e.target.value})}
            required
          />
          <input
            className={styles.input}
            type="password"
            placeholder="Password"
            value={form.password}
            onChange={e => setForm({...form, password: e.target.value})}
            required
          />
          <button className={styles.btn} disabled={loading}>
            {loading ? 'Logging in...' : 'Login'}
          </button>
        </form>
        <p className={styles.switch}>
          Don't have an account? <Link to="/register">Register</Link>
        </p>
      </div>
    </div>
  );
}