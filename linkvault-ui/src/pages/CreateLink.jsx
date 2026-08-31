import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import api, { API_URL } from '../api';
import Navbar from '../components/Navbar';
import styles from './CreateLink.module.css';

const CATEGORIES = ['General', 'Portfolio', 'Study Material', 'Projects', 'Social Media', 'Personal'];

export default function CreateLink() {
  const [form, setForm] = useState({
    original: '', alias: '', password: '', category: 'General', expires_in: ''
  });
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(null);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    const payload = { original: form.original };
    if (form.alias) payload.alias = form.alias;
    if (form.password) payload.password = form.password;
    if (form.category) payload.category = form.category;
    if (form.expires_in) payload.expires_in = form.expires_in;
    try {
      const res = await api.post('/shorten', payload);
      setResult(res.data);
      toast.success('Link created!');
    } catch (err) {
      toast.error(err.response?.data?.error || 'Failed to create link');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <Navbar />
      <div className={styles.container}>
        <div className={styles.card}>
          <h2 className={styles.title}>Create New Link</h2>
          <form onSubmit={handleSubmit} className={styles.form}>
            <div className={styles.field}>
              <label>Original URL *</label>
              <input
                className={styles.input}
                type="url"
                placeholder="https://example.com"
                value={form.original}
                onChange={e => setForm({...form, original: e.target.value})}
                required
              />
            </div>
            <div className={styles.row}>
              <div className={styles.field}>
                <label>Custom Alias</label>
                <input
                  className={styles.input}
                  placeholder="e.g. my-portfolio"
                  value={form.alias}
                  onChange={e => setForm({...form, alias: e.target.value})}
                />
              </div>
              <div className={styles.field}>
                <label>Password (optional)</label>
                <input
                  className={styles.input}
                  type="password"
                  placeholder="Protect this link"
                  value={form.password}
                  onChange={e => setForm({...form, password: e.target.value})}
                />
              </div>
            </div>
            <div className={styles.row}>
              <div className={styles.field}>
                <label>Category</label>
                <select
                  className={styles.input}
                  value={form.category}
                  onChange={e => setForm({...form, category: e.target.value})}
                >
                  {CATEGORIES.map(c => <option key={c}>{c}</option>)}
                </select>
              </div>
              <div className={styles.field}>
                <label>Expires In</label>
                <select
                  className={styles.input}
                  value={form.expires_in}
                  onChange={e => setForm({...form, expires_in: e.target.value})}
                >
                  <option value="">Never</option>
                  <option value="1h">1 Hour</option>
                  <option value="24h">24 Hours</option>
                  <option value="7d">7 Days</option>
                  <option value="30d">30 Days</option>
                </select>
              </div>
            </div>
            <button className={styles.btn} disabled={loading}>
              {loading ? 'Creating...' : 'Create Link'}
            </button>
          </form>
          {result && (
            <div className={styles.result}>
              <p className={styles.resultLabel}>Your short link is ready!</p>
              <a href={result.short_url} target="_blank" rel="noreferrer" className={styles.shortUrl}>
                {result.short_url}
              </a>
              <div className={styles.resultActions}>
                <button className={styles.copyBtn} onClick={() => {
                  navigator.clipboard.writeText(result.short_url);
                  toast.success('Copied!');
                }}>Copy</button>
                
                <a href={"http://localhost:8080/qr/" + result.short_url.split('/').pop()}
                  target="_blank"
                  rel="noreferrer"
                  className={styles.qrBtn}
                >View QR</a>
                <button className={styles.dashBtn} onClick={() => navigate('/')}>Dashboard</button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
