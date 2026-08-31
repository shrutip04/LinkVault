import { useNavigate } from 'react-router-dom';
import styles from './Navbar.module.css';

export default function Navbar() {
  const navigate = useNavigate();
  const username = localStorage.getItem('username');

  const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('username');
    navigate('/login');
  };

  return (
    <nav className={styles.nav}>
      <span className={styles.logo} onClick={() => navigate('/')}>LinkVault</span>
      <div className={styles.right}>
        <span className={styles.user}>{username}</span>
        <button className={styles.createBtn} onClick={() => navigate('/create')}>+ New Link</button>
        <button className={styles.logoutBtn} onClick={logout}>Logout</button>
      </div>
    </nav>
  );
}