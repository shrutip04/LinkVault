import { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import api, { API_URL } from '../api';
import Navbar from '../components/Navbar';
import styles from './Dashboard.module.css';

export default function Dashboard() {
  const [stats, setStats] = useState(null);
  const [links, setLinks] = useState([]);
  const [filter, setFilter] = useState('');
  const [loading, setLoading] = useState(true);

  const fetchData = async () => {
    try {
      const [statsRes, linksRes] = await Promise.all([
        api.get('/links/stats'),
        api.get('/links' + (filter ? '?category=' + filter : ''))
      ]);
      setStats(statsRes.data);
      setLinks(linksRes.data);
    } catch {
      toast.error('Failed to load dashboard');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { fetchData(); }, [filter]);

  const deleteLink = async (id) => {
    if (!window.confirm('Delete this link?')) return;
    try {
      await api.delete('/links/' + id);
      toast.success('Link deleted');
      fetchData();
    } catch {
      toast.error('Failed to delete');
    }
  };

  const copyLink = (short) => {
    navigator.clipboard.writeText(`${API_URL}/${short}`);
    toast.success('Copied!');
  };

  if (loading) return (
    <div>
      <Navbar />
      <div className={styles.loading}>Loading your dashboard...</div>
    </div>
  );

  const categoryChartData = stats && stats.categories
    ? Object.entries(stats.categories).map(function(entry) {
        return { name: entry[0], count: entry[1] };
      })
    : [];

  return (
    <div>
      <Navbar />
      <div className={styles.container}>
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <p className={styles.statLabel}>Total Links</p>
            <p className={styles.statValue}>{stats ? stats.total_links : 0}</p>
          </div>
          <div className={styles.statCard}>
            <p className={styles.statLabel}>Total Clicks</p>
            <p className={styles.statValue}>{stats ? stats.total_clicks : 0}</p>
          </div>
          <div className={styles.statCard}>
            <p className={styles.statLabel}>Active Links</p>
            <p className={styles.statValue + ' ' + styles.green}>{stats ? stats.active_links : 0}</p>
          </div>
          <div className={styles.statCard}>
            <p className={styles.statLabel}>Expired Links</p>
            <p className={styles.statValue + ' ' + styles.red}>{stats ? stats.expired_links : 0}</p>
          </div>
        </div>

        <div className={styles.bottomGrid}>
          <div className={styles.tableSection}>
            <div className={styles.tableHeader}>
              <h3>Your Links</h3>
              <select
                className={styles.filterSelect}
                value={filter}
                onChange={e => setFilter(e.target.value)}
              >
                <option value="">All Categories</option>
                {categoryChartData.map(c => (
                  <option key={c.name}>{c.name}</option>
                ))}
              </select>
            </div>
            {links.length === 0 ? (
              <p className={styles.empty}>No links yet. Create your first one!</p>
            ) : (
              <div className={styles.tableWrapper}>
                <table className={styles.table}>
                  <thead>
                    <tr>
                      <th>Short URL</th>
                      <th>Original</th>
                      <th>Clicks</th>
                      <th>Category</th>
                      <th>Status</th>
                      <th>Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {links.map(link => (
                      <tr key={link.id}>
                        <td>
                          <span className={styles.shortCode}>
                            {link.short}{link.is_protected ? ' Lock' : ''}
                          </span>
                        </td>
                        <td>
                          
                           <a href={link.original}
                            target="_blank"
                            rel="noreferrer"
                            className={styles.originalUrl}
                          >
                            {link.original.length > 35 ? link.original.slice(0, 35) + '...' : link.original}
                          </a>
                        </td>
                        <td className={styles.clicks}>{link.clicks}</td>
                        <td><span className={styles.category}>{link.category}</span></td>
                        <td>
                          <span className={link.status === 'active' ? styles.active : styles.expired}>
                            {link.status === 'active' ? 'Active' : 'Expired'}
                          </span>
                        </td>
                        <td>
                          <div className={styles.actions}>
                            <button className={styles.copyBtn} onClick={() => copyLink(link.short)}>Copy</button>
                            <a
                              href={`${API_URL}/qr/${link.short}`}
                              target="_blank"
                              rel="noreferrer"
                              className={styles.qrBtn}
                            >QR</a>
                            <button className={styles.deleteBtn} onClick={() => deleteLink(link.id)}>Del</button>
                          </div>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          {categoryChartData.length > 0 && (
            <div className={styles.chartSection}>
              <h3>Links by Category</h3>
              <ResponsiveContainer width="100%" height={250}>
                <BarChart data={categoryChartData}>
                  <XAxis dataKey="name" tick={{ fill: '#94a3b8', fontSize: 11 }} />
                  <YAxis tick={{ fill: '#94a3b8', fontSize: 11 }} />
                  <Tooltip
                    contentStyle={{ background: '#1e293b', border: '1px solid #334155', borderRadius: 8 }}
                    labelStyle={{ color: '#e2e8f0' }}
                  />
                  <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                    {categoryChartData.map((entry, i) => (
                      <Cell key={i} fill={['#6366f1','#22c55e','#f59e0b','#ec4899','#14b8a6'][i % 5]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
              {stats && stats.most_popular && stats.most_popular.short && (
                <div className={styles.popularCard}>
                  <p className={styles.popularLabel}>Most Popular</p>
                  <p className={styles.popularShort}>{stats.most_popular.short}</p>
                  <p className={styles.popularClicks}>{stats.most_popular.clicks} clicks</p>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
