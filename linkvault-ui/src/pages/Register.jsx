import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import api from '../api';
import styles from './Auth.module.css';

export default function Register() {
  const [form, setForm] = useState({ username: '', email: '', password: '' });
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      await api.post('/register', form);
      toast.success('Account created! Please login.');
      navigate('/login');
    } catch (err) {
      toast.error(err.response?.data?.error || 'Registration failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <h1 className={styles.logo}>🔗 LinkVault</h1>
        <h2 className={styles.title}>Create account</h2>
        <form onSubmit={handleSubmit} className={styles.form}>
          <input
            className={styles.input}
            type="text"
            placeholder="Username"
            value={form.username}
            onChange={e => setForm({...form, username: e.target.value})}
            required
          />
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
            {loading ? 'Creating...' : 'Register'}
          </button>
        </form>
        <p className={styles.switch}>
          Already have an account? <Link to="/login">Login</Link>
        </p>
      </div>
    </div>
  );
}